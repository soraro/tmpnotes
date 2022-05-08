package config

import (
	"html/template"
)

var (
	Tmpl *template.Template
	err  error
	path = "templates/*"
)

func GetTemplates() error {
	Tmpl, err = template.ParseGlob(path)
	if err != nil {
		return err
	}
	return nil
}
