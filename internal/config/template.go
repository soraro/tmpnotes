package config

import (
	"html/template"
)

var (
	Tmpl *template.Template
	err  error
)

func GetTemplates() error {
	Tmpl, err = template.ParseGlob("templates/*")
	if err != nil {
		return err
	}
	return nil
}
