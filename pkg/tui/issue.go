package tui

import (
	"bytes"
	"text/template"
)

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

func (d *DashboardTab) renderTemplate(issue *IssueTemplate, templateFile string) (output bytes.Buffer, err error) {
	var file *template.Template
	if file, err = template.ParseFiles(templateFile); err != nil {
		return output, err
	}
	if err = file.Execute(&output, issue); err != nil {
		return output, err
	}
	return
}
