package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nndi-oss/greypot"
	"github.com/nndi-oss/greypot/http/fiber/handlers"
)

func Use(app fiber.Router, s *greypot.Module) {
	app.Get("/reports/list", handlers.ReportListHandler(s.TemplateService))
	app.Get("/reports/preview/*", handlers.ReportPreviewHandler(s.TemplateService, s.TemplateEngine, s.ReportService))
	app.Get("/reports/render/*", handlers.ReportRenderHandlder(s.TemplateService))

	app.Post("/reports/export/html/*", handlers.ReportExportHandler(s.ReportService, "html"))
	app.Post("/reports/export/png/*", handlers.ReportExportHandler(s.ReportService, "png"))
	app.Post("/reports/export/pdf/*", handlers.ReportExportHandler(s.ReportService, "pdf"))

	app.Post("/reports/export/bulk/html/*", handlers.BulkReportExportHandler(s.ReportService, "html"))
	app.Post("/reports/export/bulk/png/*", handlers.BulkReportExportHandler(s.ReportService, "png"))
	app.Post("/reports/export/bulk/pdf/*", handlers.BulkReportExportHandler(s.ReportService, "pdf"))
}
