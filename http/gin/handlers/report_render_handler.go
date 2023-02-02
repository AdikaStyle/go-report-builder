package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nndi-oss/greypot/service"
	"github.com/sirupsen/logrus"
)

func ReportRenderHandlder(templateService service.TemplateService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reportId := strings.TrimPrefix(ctx.Param("reportId"), "/")
		data, err := extractData(ctx)
		if err != nil {
			logrus.Error(err.Error())
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		html, err := templateService.RenderTemplate(reportId, data)
		if err != nil {
			logrus.Error(err.Error())
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.Data(http.StatusOK, "text/html; charset=utf-8", html)
	}
}
