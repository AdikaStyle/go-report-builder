package greypot

import (
	"fmt"
	"io/fs"
	"time"

	"github.com/nndi-oss/greypot/exporter"
	"github.com/nndi-oss/greypot/service"
	"github.com/nndi-oss/greypot/template/engine"
	"github.com/nndi-oss/greypot/template/repo"
	"github.com/playwright-community/playwright-go"
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

func NewChromeDPModule(renderTimeout time.Duration, repo repo.TemplateRepository) *Module {
	return NewModule(
		WithRenderTimeout(renderTimeout),
		WithViewport(2048, 1920),
		WithTemplateEngine(engine.NewGolangTemplateEngine()),
		WithTemplatesRepository(repo),
		WithChromeDPRenderer("localhost"),
	)
}

func NewPlaywrightModule(renderTimeout time.Duration, repo repo.TemplateRepository) *Module {
	err := playwright.Install()
	if err != nil {
		panic(fmt.Errorf("could not install playwright dependencies (Chromium): %w", err))
	}

	return NewModule(
		WithRenderTimeout(renderTimeout),
		WithViewport(2048, 1920),
		WithTemplateEngine(engine.NewGolangTemplateEngine()),
		WithTemplatesRepository(repo),
		WithChromeDPRenderer("localhost"),
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
