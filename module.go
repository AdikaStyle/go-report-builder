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
	TemplateService service.TemplateService
	ReportService   service.ReportService

	TemplateEngine     engine.TemplateEngine
	TemplateRepository repo.TemplateRepository
	PDFExporter        exporter.ReportExporter
	PNGExporter        exporter.ReportExporter
}

func NewModule(renderTimeout time.Duration, repo repo.TemplateRepository) *Module {
	var realRenderTimeout time.Duration
	if renderTimeout == 0 {
		realRenderTimeout = 10 * time.Second
	} else {
		realRenderTimeout = renderTimeout
	}
	var viewportHeight int = 2048 // `env:"VIEWPORT_HEIGHT" envDefault:"2048"`
	var viewportWidth int = 1920  // `env:"VIEWPORT_WIDTH" envDefault:"1920"`

	engine := engine.NewGolangTemplateEngine()
	// repo := data.NewFilesystemTemplateRepo(config.TemplatesPath)
	png := exporter.NewChromePNGReportExporter(realRenderTimeout, viewportHeight, viewportWidth)
	pdf := exporter.NewChromePDFReportExporter(realRenderTimeout, viewportHeight, viewportWidth)

	tmplSrv := service.NewTemplateService(engine, repo)
	rprtSrv := service.NewReportService(pdf, png, tmplSrv, "/" /* TODO: remove baseURL param?*/)

	return &Module{
		TemplateService:    tmplSrv,
		ReportService:      rprtSrv,
		TemplateEngine:     engine,
		TemplateRepository: repo,
		PDFExporter:        pdf,
		PNGExporter:        png,
	}
}

func NewPlaywrightModule(renderTimeout time.Duration, repo repo.TemplateRepository) *Module {
	err := playwright.Install()
	if err != nil {
		panic(fmt.Errorf("could not install playwright dependencies (Chromium): %w", err))
	}

	var realRenderTimeout time.Duration
	if renderTimeout == 0 {
		realRenderTimeout = 10 * time.Second
	} else {
		realRenderTimeout = renderTimeout
	}
	var viewportHeight int = 2048 // `env:"VIEWPORT_HEIGHT" envDefault:"2048"`
	var viewportWidth int = 1920  // `env:"VIEWPORT_WIDTH" envDefault:"1920"`

	engine := engine.NewGolangTemplateEngine()
	// repo := data.NewFilesystemTemplateRepo(config.TemplatesPath)
	png := exporter.NewPlaywrightPNGReportExporter(realRenderTimeout, viewportHeight, viewportWidth)
	pdf := exporter.NewPlaywrightPDFReportExporter(realRenderTimeout, viewportHeight, viewportWidth)

	tmplSrv := service.NewTemplateService(engine, repo)
	rprtSrv := service.NewReportService(pdf, png, tmplSrv, "/" /* TODO: remove baseURL param?*/)

	return &Module{
		TemplateService:    tmplSrv,
		ReportService:      rprtSrv,
		TemplateEngine:     engine,
		TemplateRepository: repo,
		PDFExporter:        pdf,
		PNGExporter:        png,
	}
}

func NewPlaywrightModuleWithCustomEngine(renderTimeout time.Duration, repo repo.TemplateRepository, templateEngine engine.TemplateEngine) *Module {
	err := playwright.Install()
	if err != nil {
		panic(fmt.Errorf("could not install playwright dependencies (Chromium): %w", err))
	}

	var realRenderTimeout time.Duration
	if renderTimeout == 0 {
		realRenderTimeout = 10 * time.Second
	} else {
		realRenderTimeout = renderTimeout
	}
	var viewportHeight int = 2048 // `env:"VIEWPORT_HEIGHT" envDefault:"2048"`
	var viewportWidth int = 1920  // `env:"VIEWPORT_WIDTH" envDefault:"1920"`

	// repo := data.NewFilesystemTemplateRepo(config.TemplatesPath)
	png := exporter.NewPlaywrightPNGReportExporter(realRenderTimeout, viewportHeight, viewportWidth)
	pdf := exporter.NewPlaywrightPDFReportExporter(realRenderTimeout, viewportHeight, viewportWidth)

	tmplSrv := service.NewTemplateService(templateEngine, repo)
	rprtSrv := service.NewReportService(pdf, png, tmplSrv, "/" /* TODO: remove baseURL param?*/)

	return &Module{
		TemplateService:    tmplSrv,
		ReportService:      rprtSrv,
		TemplateEngine:     templateEngine,
		TemplateRepository: repo,
		PDFExporter:        pdf,
		PNGExporter:        png,
	}
}

func NewFSTemplateRepo(templatesFolder fs.FS) repo.TemplateRepository {
	return repo.NewFSTemplateRepo(templatesFolder)
}

func NewFilesystemRepo(templatesFolder string) repo.TemplateRepository {
	return repo.NewFilesystemTemplateRepo(templatesFolder)
}
