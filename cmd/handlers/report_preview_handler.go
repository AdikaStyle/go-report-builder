package handlers

import (
	"github.com/AdikaStyle/go-report-builder/internal/business"
	"github.com/AdikaStyle/go-report-builder/internal/data"
	"github.com/AdikaStyle/go-report-builder/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func ReportPreviewHandler(templateEngine data.TemplateEngine, reportService business.ReportService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reportId := ctx.Param("reportId")
		data, err := extractData(ctx)
		if err != nil {
			logrus.Error(err)
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		png, err := reportService.ExportReportPng(reportId, data)
		if err != nil {
			logrus.Error(err)
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		rendered, err := templateEngine.Render(
			[]byte(previewTemplate),
			&models.TemplateContext{Values: map[string]interface{}{
				"reportId": reportId,
				"data":     data,
				"image":    string(png),
			}},
		)
		if err != nil {
			logrus.Error(err)
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.Data(http.StatusOK, "text/html; charset=utf-8", rendered)
	}
}

const previewTemplate = `
<!DOCTYPE html>
<html>
<head>
	<style>
		
	</style>
</head>
<body>

<div>
	Preview App:
	<p> {{ .Values.reportId }} </p>
	<p> {{ .Values.data }} </p>
	<img src="data:image/png;base64, {{ .Values.image }}" alt="Report" />
</div>

</body>
</html>
`
