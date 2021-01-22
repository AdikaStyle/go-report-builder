package cmd

import (
	"github.com/AdikaStyle/go-report-builder/internal/business"
	"github.com/AdikaStyle/go-report-builder/internal/data"
)

type Module struct {
	Config Config

	templateService business.TemplateService
	reportService   business.ReportService

	templateEngine     data.TemplateEngine
	templateRepository data.TemplateRepository
	pdfExporter        data.ReportExporter
	pngExporter        data.ReportExporter

	Server *Server
}

func NewModule(config Config) *Module {
	engine := data.NewGolangTemplateEngine()
	repo := data.NewFilesystemTemplateRepo(config.TemplatesPath)
	png := data.NewPngReportExporter(config.RenderTimeout)
	pdf := data.NewPdfReportExporter(png)

	tmplSrv := business.NewTemplateService(engine, repo)
	rprtSrv := business.NewReportService(pdf, png, tmplSrv, config.BaseUrl())

	srv := NewServer(config, tmplSrv, engine, rprtSrv)

	return &Module{
		Config:             config,
		templateService:    tmplSrv,
		reportService:      rprtSrv,
		templateEngine:     engine,
		templateRepository: repo,
		pdfExporter:        pdf,
		pngExporter:        png,
		Server:             srv,
	}
}
