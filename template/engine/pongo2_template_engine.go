package engine

import (
	"bytes"

	"github.com/flosch/pongo2/v6"
	"github.com/nndi-oss/greypot/models"
	"github.com/palantir/stacktrace"
)

type pongo2TemplateEngine struct{}

func NewDjangoTemplateEngine() *pongo2TemplateEngine {
	return &pongo2TemplateEngine{}
}

func NewPongo2TemplateEngine() *pongo2TemplateEngine {
	return &pongo2TemplateEngine{}
}

func (pte *pongo2TemplateEngine) Render(templateContent []byte, ctx *models.TemplateContext) ([]byte, error) {
	t, err := pongo2.FromBytes(templateContent)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to parse template file: %s with data: %+v", templateContent, ctx)
	}

	pongoCtx := pongo2.Context{"data": ctx.Values}
	var out bytes.Buffer
	if err := t.ExecuteWriter(pongoCtx, &out); err != nil {
		return nil, stacktrace.Propagate(err, "failed to execute template file: %s with data: %+v", templateContent, ctx)
	}

	return out.Bytes(), err
}
