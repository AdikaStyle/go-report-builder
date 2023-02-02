package handlers

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/nndi-oss/greypot/service"
	"github.com/sirupsen/logrus"
)

func ReportRenderHandlder(templateService service.TemplateService) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		reportId := strings.TrimPrefix(ctx.Params("*"), "/")
		data, err := extractData(ctx)
		if err != nil {
			logrus.Error(err.Error())
			return ctx.Status(http.StatusInternalServerError).
				JSON(fiber.Map{
					"err": err.Error(),
				})
		}

		html, err := templateService.RenderTemplate(reportId, data)
		if err != nil {
			logrus.Error(err.Error())
			return ctx.Status(http.StatusInternalServerError).
				JSON(fiber.Map{
					"err": err.Error(),
				})
		}

		ctx.Type(".html", "charset=utf-8")
		return ctx.Send(html)
	}
}
