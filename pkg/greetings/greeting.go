package greetings

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"
)

type Greeting struct {
	tplBody string
	tpl     *template.Template
}

func NewGreeting(tplBody string) (*Greeting, error) {
	tpl, err := template.New("Greeting").Parse(tplBody)
	if err != nil {
		return nil, fmt.Errorf("failed to parse greeting template body: %w", err)
	}

	return &Greeting{
		tplBody: tplBody,
		tpl:     tpl,
	}, nil
}

func (g *Greeting) Greet(name string) (string, error) {
	if name == "" {
		return "", errors.New("name should not be empty")
	}

	var buff bytes.Buffer
	if err := g.tpl.Execute(&buff, name); err != nil {
		return "", fmt.Errorf("failed to execute greeting template: %w", err)
	}

	return buff.String(), nil
}

func (g *Greeting) Template() string {
	return g.tplBody
}
