package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nndi-oss/greypot/service"
	"github.com/sirupsen/logrus"
)

func ReportExportHandler(reportService service.ReportService, kind string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reportId := strings.TrimPrefix(ctx.Param("reportId"), "/")
		var body interface{}
		if err := ctx.BindJSON(&body); err != nil {
			logrus.Error(err)
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		var export []byte
		var err error
		switch kind {
		case "html":
			export, err = reportService.ExportReportHtml(reportId, body)
		case "pdf":
			export, err = reportService.ExportReportPdf(reportId, body)
		case "png":
			export, err = reportService.ExportReportPng(reportId, body)
		}

		if err != nil {
			logrus.Error(err)
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, map[string]interface{}{
			"reportId": reportId,
			"data":     string(export),
			"type":     kind,
		})
	}
}
