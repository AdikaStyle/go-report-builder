package cmd

import (
	"fmt"
	"github.com/AdikaStyle/go-report-builder/cmd/handlers"
	"github.com/AdikaStyle/go-report-builder/internal/business"
	"github.com/AdikaStyle/go-report-builder/internal/data"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config          Config
	templateService business.TemplateService
	reportService   business.ReportService
	templateEngine  data.TemplateEngine
}

func NewServer(config Config, templateService business.TemplateService, engine data.TemplateEngine, reportService business.ReportService) *Server {
	return &Server{config: config, templateService: templateService, reportService: reportService, templateEngine: engine}
}

func (s *Server) setupServer() *gin.Engine {
	r := gin.Default()
	r.GET("/health", handlers.Health())
	r.GET("/reports/list", handlers.ReportListHandler(s.templateService))
	r.GET("/reports/render/*reportId", handlers.ReportRenderHandlder(s.templateService))
	r.GET("/reports/preview/*reportId", handlers.ReportPreviewHandler(s.templateService, s.templateEngine, s.reportService))
	r.POST("/reports/export/html/*reportId", handlers.ReportExportHandler(s.reportService, "html"))
	r.POST("/reports/export/png/*reportId", handlers.ReportExportHandler(s.reportService, "png"))
	r.POST("/reports/export/pdf/*reportId", handlers.ReportExportHandler(s.reportService, "pdf"))
	return r
}

func (s *Server) Start() error {
	srv := s.setupServer()
	return srv.Run(fmt.Sprintf(":%d", s.config.ServerPort))
}
