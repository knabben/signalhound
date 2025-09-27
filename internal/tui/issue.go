package tui

import (
	"bytes"
	"embed"
	"text/template"
)

//go:embed template/*
var tmplFolder embed.FS

type IssueTemplate struct {
	BoardName    string
	TabName      string
	TestName     string
	FirstFailure string
	LastFailure  string
	TestGridURL  string
	TriageURL    string
	ProwURL      string
	ErrMessage   string
	Sig          string
}

func renderTemplate(issue *IssueTemplate, templateFile string) (output bytes.Buffer, err error) {
	var tmpl *template.Template
	tmpl, err = template.ParseFS(tmplFolder, templateFile)
	if err != nil {
		return output, err
	}
	if err = tmpl.Execute(&output, issue); err != nil {
		return output, err
	}
	return
}
