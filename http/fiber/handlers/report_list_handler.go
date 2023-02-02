package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/nndi-oss/greypot/service"
	"github.com/sirupsen/logrus"
)

func ReportListHandler(tmplSrv service.TemplateService) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		templates, err := tmplSrv.ListTemplates()
		if err != nil {
			logrus.Error(err.Error())
			return ctx.Status(http.StatusInternalServerError).
				JSON(fiber.Map{
					"err": err.Error(),
				})
		}

		return ctx.JSON(fiber.Map{
			"list": templates,
		})
	}
}
