package notes

import (
	"bytes"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"

	cfg "tmpnotes/internal/config"
)

// LOOK AT THIS: https://github.com/slack-go/slack/blob/master/examples/modal/modal.go

var api *slack.Client
var SlackEnabled bool = false

func SlackInit() {
	api = slack.New(cfg.Config.SlackToken)
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

	modalRequest := generateModalRequest()
	_, err = api.OpenView(s.TriggerID, modalRequest)
	if err != nil {
		log.Errorf("Error opening view: %s\n", err)
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
	note := slack.NewInputBlock("Note", noteText, nil, noteElement)

	var opts []*slack.OptionBlockObject
	expireText := slack.NewTextBlockObject("plain_text", "Hours before disappearing", false, false)
	opts = append(opts, slack.NewOptionBlockObject("value-1", slack.NewTextBlockObject("plain_text", "1", false, false), nil))
	opts = append(opts, slack.NewOptionBlockObject("value-2", slack.NewTextBlockObject("plain_text", "2", false, false), nil))
	opts = append(opts, slack.NewOptionBlockObject("value-3", slack.NewTextBlockObject("plain_text", "3", false, false), nil))
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
