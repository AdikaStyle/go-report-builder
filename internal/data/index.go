package data

import "github.com/AdikaStyle/go-report-builder/internal/models"

type TemplateEngine interface {
	Render(templateContent []byte, ctx *models.TemplateContext) ([]byte, error)
}

type TemplateRepository interface {
	ListAll() ([]string, error)
	LoadTemplate(templateId string) ([]byte, error)
}

type ReportExporter interface {
	Export(renderedContent []byte) ([]byte, error)
}
