package notes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"

	cfg "tmpnotes/internal/config"
)

// LOOK AT THIS: https://github.com/slack-go/slack/blob/master/examples/modal/modal.go

var api *slack.Client
var modal slack.ModalViewRequest
var SlackEnabled bool = false

func SlackInit() {
	api = slack.New(cfg.Config.SlackToken)
	modal = generateModalRequest()
	SlackEnabled = true
}

func verifySigningSecret(r *http.Request) error {
	verifier, err := slack.NewSecretsVerifier(r.Header, cfg.Config.SlackSigningSecret)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	// Need to use r.Body again when unmarshalling SlashCommand and InteractionCallback
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	verifier.Write(body)
	if err = verifier.Ensure(); err != nil {
		return err
	}

	return nil
}

func SlackHandler(w http.ResponseWriter, r *http.Request) {
	err := verifySigningSecret(r)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	s, err := slack.SlashCommandParse(r)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = api.OpenView(s.TriggerID, modal)
	if err != nil {
		log.Errorf("Error opening view: %s\n", err)
	}

}

func SlackResponseHandler(w http.ResponseWriter, r *http.Request) {
	err := verifySigningSecret(r)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var i slack.InteractionCallback
	err = json.Unmarshal([]byte(r.FormValue("payload")), &i)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
	}

	note := i.View.State.Values["Note"]["note"].Value
	expire := i.View.State.Values["Expire"]["Expire"].SelectedOption.Text.Text

	fmt.Println(note, expire)

	id, key := generateIdAndKey()

	encryptedMessage, err := encryptNote(note, key)
	if err != nil {
		log.Errorf("%s Issue encrypting message: %s", r.RequestURI, err)
	}
	pipe := rdb.Pipeline()
	expireDuration, _ := strconv.Atoi(expire)
	pipe.Set(ctx, id, encryptedMessage, time.Duration(expireDuration)*time.Hour)
	pipe.HIncrBy(ctx, "counts", noteType(note), 1)
	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Errorf("%s Error setting note values: %s", r.RequestURI, err)
		http.Error(w, "Error connecting to database", 500)
		return
	}

	_, _, err = api.PostMessage(
		i.User.ID,
		slack.MsgOptionText(fmt.Sprintf("Your tmpnote: %s://%s/id/%s%s", r.Header["X-Forwarded-Proto"][0], r.Host, id, key), false),
		slack.MsgOptionAttachments(),
	)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
	}

}

func generateModalRequest() slack.ModalViewRequest {
	titleText := slack.NewTextBlockObject("plain_text", "TMPNOTES 🔥️", false, false)
	closeText := slack.NewTextBlockObject("plain_text", "Close", false, false)
	submitText := slack.NewTextBlockObject("plain_text", "Submit", false, false)

	headerText := slack.NewTextBlockObject("plain_text", "Enter your secret to create a link in this channel", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	noteText := slack.NewTextBlockObject("plain_text", "Note:", false, false)
	notePlaceholder := slack.NewTextBlockObject("plain_text", "Write something", false, false)
	noteElement := slack.NewPlainTextInputBlockElement(notePlaceholder, "note")
	noteElement.Multiline = true
	noteElement.MaxLength = cfg.Config.UiMaxLength
	note := slack.NewInputBlock("Note", noteText, nil, noteElement)

	// https://api.slack.com/reference/block-kit/block-elements#conversations_select
	//channelDialog := slack.DialogInputSelect{}
	//channelText := slack.NewTextBlockObject("plain_text", "Select a channel from the list", false, false)
	//channelDialog := slack.NewConversationsSelect("Channel Select", "select a channel")

	var opts []*slack.OptionBlockObject
	expireText := slack.NewTextBlockObject("plain_text", "Hours before disappearing", false, false)
	for i := 1; i <= cfg.Config.MaxExpire; i++ {
		opts = append(opts, slack.NewOptionBlockObject(fmt.Sprintf("value-%v", i), slack.NewTextBlockObject("plain_text", fmt.Sprint(i), false, false), nil))
	}
	optionGroup := slack.NewOptionGroupBlockElement(expireText, opts...)
	expireElement := slack.NewOptionsGroupSelectBlockElement("static_select", nil, "Expire", optionGroup)
	expire := slack.NewInputBlock("Expire", expireText, nil, expireElement)

	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			headerSection,
			note,
			expire,
		},
	}

	var modalRequest slack.ModalViewRequest
	modalRequest.Type = slack.ViewType("modal")
	modalRequest.Title = titleText
	modalRequest.Close = closeText
	modalRequest.Submit = submitText
	modalRequest.Blocks = blocks
	return modalRequest
}
