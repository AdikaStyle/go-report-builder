package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/nndi-oss/greypot"
	"github.com/nndi-oss/greypot/http/gin/handlers"
)

func Use(r gin.IRouter, s *greypot.Module) {
	r.GET("/health", handlers.Health())
	r.GET("/reports/list", handlers.ReportListHandler(s.TemplateService))
	r.GET("/reports/render/*reportId", handlers.ReportRenderHandlder(s.TemplateService))
	r.GET("/reports/preview/*reportId", handlers.ReportPreviewHandler(s.TemplateService, s.TemplateEngine, s.ReportService))

	r.POST("/reports/export/html/*reportId", handlers.ReportExportHandler(s.ReportService, "html"))
	r.POST("/reports/export/png/*reportId", handlers.ReportExportHandler(s.ReportService, "png"))
	r.POST("/reports/export/pdf/*reportId", handlers.ReportExportHandler(s.ReportService, "pdf"))

	r.POST("/reports/export/bulk/html/*reportId", handlers.BulkReportExportHandler(s.ReportService, "html"))
	r.POST("/reports/export/bulk/png/*reportId", handlers.BulkReportExportHandler(s.ReportService, "png"))
	r.POST("/reports/export/bulk/pdf/*reportId", handlers.BulkReportExportHandler(s.ReportService, "pdf"))
}
