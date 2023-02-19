package handlers

import (
	"fmt"
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

type ExportResponse struct {
	ID   string `json:"reportId"`
	Data string `json:"data"`
	Type string `json:"type"`
}

type BulkExportEntry struct {
	ID   string `json:"_id"`
	Data any    `json:"data"`
}

type BulkExportRequest struct {
	ID      string            `json:"_id"`
	Entries []BulkExportEntry `json:"entries"`
}

type BulkExportResponse struct {
	ID       string           `json:"_id"`
	ReportID string           `json:"reportId"`
	Reports  []ExportResponse `json:"reports"`
}

func BulkReportExportHandler(reportService service.ReportService, kind string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reportId := strings.TrimPrefix(ctx.Param("reportId"), "/")
		body := new(BulkExportRequest)
		if err := ctx.BindJSON(&body); err != nil {
			logrus.Error(err)
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		bulkResponse := &BulkExportResponse{
			ID:       body.ID,
			ReportID: reportId,
			Reports:  make([]ExportResponse, 0),
		}

		for _, entry := range body.Entries {
			entryReal := entry
			reportData := entryReal.Data
			var export []byte
			var err error
			switch kind {
			case "html":
				export, err = reportService.ExportReportHtml(reportId, reportData)
			case "pdf":
				export, err = reportService.ExportReportPdf(reportId, reportData)
			case "png":
				export, err = reportService.ExportReportPng(reportId, reportData)
			}
			if err != nil {
				logrus.Error(err)
				ctx.String(http.StatusInternalServerError, err.Error())
				return
			}

			bulkResponse.Reports = append(bulkResponse.Reports, ExportResponse{
				ID:   fmt.Sprintf("%s:%s", reportId, entryReal.ID),
				Data: string(export),
				Type: kind,
			})
		}

		ctx.JSON(http.StatusOK, bulkResponse)
	}
}
