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

type GreypotOption func(*Module)

func NewModule(options ...GreypotOption) *Module {
	mod := moduleDefaults()

	for _, co := range options {
		if co == nil {
			continue
		}
		co(mod)
	}

	mod.TemplateService = service.NewTemplateService(mod.TemplateEngine, mod.TemplateRepository)
	mod.ReportService = service.NewReportService(mod.PDFExporter, mod.PNGExporter, mod.TemplateService, "/" /* TODO: remove baseURL param?*/)

	return mod
}

func WithRenderTimeout(timeout time.Duration) GreypotOption {
	return func(m *Module) {
		m.RenderTimeout = timeout
	}
}

func WithViewport(width, height int) GreypotOption {
	return func(m *Module) {
		m.ViewportWidth = width
		m.ViewportHeight = height
	}
}

func WithTemplatesRepository(repository repo.TemplateRepository) GreypotOption {
	return func(m *Module) {
		m.TemplateRepository = repository
	}
}

func WithTemplateEngine(e engine.TemplateEngine) GreypotOption {
	return func(m *Module) {
		m.TemplateEngine = e
	}
}

func WithGolangTemplateEngine() GreypotOption {
	return func(m *Module) {
		m.TemplateEngine = engine.NewGolangTemplateEngine()
	}
}

func WithDjangoTemplateEngine() GreypotOption {
	return func(m *Module) {
		m.TemplateEngine = engine.NewDjangoTemplateEngine()
	}
}

func WithPlaywrightRenderer(options ...*playwright.RunOptions) GreypotOption {

	return func(m *Module) {
		err := playwright.Install(options...)
		if err != nil {
			panic(fmt.Errorf("could not install playwright dependencies (Chromium): %w", err))
		}
		png := exporter.NewPlaywrightPNGReportExporter(m.RenderTimeout, m.ViewportHeight, m.ViewportWidth)
		pdf := exporter.NewPlaywrightPDFReportExporter(m.RenderTimeout, m.ViewportHeight, m.ViewportWidth)
		m.PNGExporter = png
		m.PDFExporter = pdf
	}
}

func WithChromeDPRenderer(url string) GreypotOption {
	return func(m *Module) {
		png := exporter.NewChromePNGReportExporter(m.RenderTimeout, m.ViewportHeight, m.ViewportWidth)
		pdf := exporter.NewChromePDFReportExporter(m.RenderTimeout, m.ViewportHeight, m.ViewportWidth)
		m.PNGExporter = png
		m.PDFExporter = pdf
	}
}

func WithTemplatesFromFS(dir fs.FS) GreypotOption {
	return func(m *Module) {
		m.TemplateRepository = repo.NewFSTemplateRepo(dir)
	}
}

func WithTemplatesFromFilesytem(dir string) GreypotOption {
	return func(m *Module) {
		m.TemplateRepository = repo.NewFilesystemTemplateRepo(dir)
	}
}

func moduleDefaults() *Module {
	defaultRenderTimeout := 10 * time.Second

	var viewportHeight int = 2048 // `env:"VIEWPORT_HEIGHT" envDefault:"2048"`
	var viewportWidth int = 1920  // `env:"VIEWPORT_WIDTH" envDefault:"1920"`

	templateEngine := engine.NewGolangTemplateEngine()
	repo := repo.NewFilesystemTemplateRepo("./")
	// repo := data.NewFilesystemTemplateRepo(config.TemplatesPath)
	png := exporter.NewPlaywrightPNGReportExporter(defaultRenderTimeout, viewportHeight, viewportWidth)
	pdf := exporter.NewPlaywrightPDFReportExporter(defaultRenderTimeout, viewportHeight, viewportWidth)

	tmplSrv := service.NewTemplateService(templateEngine, repo)
	rprtSrv := service.NewReportService(pdf, png, tmplSrv, "/")

	return &Module{
		RenderTimeout:      defaultRenderTimeout,
		TemplateService:    tmplSrv,
		ReportService:      rprtSrv,
		TemplateEngine:     templateEngine,
		TemplateRepository: repo,
		PDFExporter:        pdf,
		PNGExporter:        png,
	}
}
