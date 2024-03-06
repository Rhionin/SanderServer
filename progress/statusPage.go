package progress

import (
	"bytes"
	_ "embed"
	"text/template"
)

//go:embed status-page.html.tmpl
var statusPageContent string

func CreateStatusPage(wips []WorkInProgress) ([]byte, error) {
	t := template.New("statusPage")

	t, err := t.Parse(statusPageContent)
	if err != nil {
		return nil, err
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, wips); err != nil {
		return nil, err
	}

	return tpl.Bytes(), nil
}
