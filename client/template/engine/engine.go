package engine

import "github.com/nndi-oss/greypot/models"

type TemplateEngine interface {
	Render(templateContent []byte, ctx *models.TemplateContext) ([]byte, error)
}
