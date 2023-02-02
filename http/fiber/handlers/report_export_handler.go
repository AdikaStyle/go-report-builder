package handlers

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/nndi-oss/greypot/service"
	"github.com/sirupsen/logrus"
)

func ReportExportHandler(reportService service.ReportService, kind string) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		reportId := strings.TrimPrefix(ctx.Params("*"), "/")
		var body interface{}
		if err := ctx.BodyParser(&body); err != nil {
			logrus.Error(err)
			return ctx.Status(http.StatusInternalServerError).
				JSON(fiber.Map{
					"err": err.Error(),
				})
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
			return ctx.Status(http.StatusInternalServerError).
				JSON(fiber.Map{
					"err": err.Error(),
				})
		}

		return ctx.JSON(fiber.Map{
			"reportId": reportId,
			"data":     string(export),
			"type":     kind,
		})
	}
}
