package data

import (
	"bytes"
	"github.com/AdikaStyle/go-report-builder/internal/models"
	"github.com/palantir/stacktrace"
	"html/template"
)

type golangTemplateEngine struct{}

func NewGolangTemplateEngine() *golangTemplateEngine {
	return &golangTemplateEngine{}
}

func (gte *golangTemplateEngine) Render(templateContent []byte, ctx *models.TemplateContext) ([]byte, error) {
	eng := template.New("template")
	t, err := eng.Parse(string(templateContent))
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to parse template file: %s with data: %+v", templateContent, ctx)
	}

	var out bytes.Buffer
	if err := t.Execute(&out, ctx); err != nil {
		return nil, stacktrace.Propagate(err, "failed to execute template file: %s with data: %+v", templateContent, ctx)
	}

	return out.Bytes(), err
}
