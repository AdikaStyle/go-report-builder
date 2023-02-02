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

func NewFSTemplateRepo(templatesFolder fs.FS) repo.TemplateRepository {
	return repo.NewFSTemplateRepo(templatesFolder)
}
