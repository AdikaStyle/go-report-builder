package greypot

import (
	"io/fs"
	"time"

	"github.com/nndi-oss/greypot/exporter"
	"github.com/nndi-oss/greypot/service"
	"github.com/nndi-oss/greypot/template/engine"
	"github.com/nndi-oss/greypot/template/repo"
)

type Module struct {
	RenderTimeout  time.Duration
	ViewportHeight int
	ViewportWidth  int

	TemplateService service.TemplateService
	ReportService   service.ReportService

	TemplateEngine     engine.TemplateEngine
	TemplateRepository repo.TemplateRepository
	PDFExporter        exporter.ReportExporter
	PNGExporter        exporter.ReportExporter
}

func NewPlaywrightModule(renderTimeout time.Duration, repo repo.TemplateRepository) *Module {
	return NewModule(
		WithRenderTimeout(renderTimeout),
		WithViewport(2048, 1920),
		WithTemplateEngine(engine.NewGolangTemplateEngine()),
		WithTemplatesRepository(repo),
		WithPlaywrightRenderer(),
	)
}

func NewFSTemplateRepo(templatesFolder fs.FS) repo.TemplateRepository {
	return repo.NewFSTemplateRepo(templatesFolder)
}

func NewFilesystemRepo(templatesFolder string) repo.TemplateRepository {
	return repo.NewFilesystemTemplateRepo(templatesFolder)
}

func NewGolangTemplateEngine() engine.TemplateEngine {
	return engine.NewGolangTemplateEngine()
}

func NewPongo2TemplateEngine() engine.TemplateEngine {
	return engine.NewDjangoTemplateEngine()
}

func NewDjangoTemplateEngine() engine.TemplateEngine {
	return engine.NewDjangoTemplateEngine()
}
