package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nndi-oss/greypot"
	"github.com/nndi-oss/greypot/http/fiber/handlers"
)

// type Server struct {
// 	config          Config
// 	templateService business.TemplateService
// 	reportService   business.ReportService
// 	templateEngine  data.TemplateEngine
// }

// func NewServer(config Config, templateService business.TemplateService, engine data.TemplateEngine, reportService business.ReportService) *Server {
// 	return &Server{config: config, templateService: templateService, reportService: reportService, templateEngine: engine}
// }s

func Use(app *fiber.App, s *greypot.Module) {
	app.Get("/reports/list", handlers.ReportListHandler(s.TemplateService))
	app.Get("/reports/preview/*", handlers.ReportPreviewHandler(s.TemplateService, s.TemplateEngine, s.ReportService))
	app.Get("/reports/render/*", handlers.ReportRenderHandlder(s.TemplateService))
	app.Post("/reports/export/html/*", handlers.ReportExportHandler(s.ReportService, "html"))
	app.Post("/reports/export/png/*", handlers.ReportExportHandler(s.ReportService, "png"))
	app.Post("/reports/export/pdf/*", handlers.ReportExportHandler(s.ReportService, "pdf"))
}
