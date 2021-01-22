package business

import (
	"github.com/AdikaStyle/go-report-builder/internal/data"
	"github.com/AdikaStyle/go-report-builder/internal/models"
	"github.com/palantir/stacktrace"
)

type templateService struct {
	engine data.TemplateEngine
	repo   data.TemplateRepository
}

func (ts *templateService) ListTemplates() ([]string, error) {
	templates, err := ts.repo.ListAll()
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to list all templates")
	}
	return templates, nil
}

func (ts *templateService) GetTemplate(reportId string) ([]byte, error) {
	content, err := ts.repo.LoadTemplate(reportId)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to get single template with id: %s", reportId)
	}

	return content, nil
}

func (ts *templateService) RenderTemplate(reportId string, data interface{}) ([]byte, error) {
	ctx := &models.TemplateContext{
		Values: data,
	}

	tmpl, err := ts.repo.LoadTemplate(reportId)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to get single template with id: %s", reportId)
	}

	html, err := ts.engine.Render(tmpl, ctx)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to render template with id: %s and data: %+v", reportId, ctx)
	}

	return html, nil
}

func NewTemplateService(engine data.TemplateEngine, repo data.TemplateRepository) *templateService {
	return &templateService{engine: engine, repo: repo}
}
