package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nndi-oss/greypot/service"
	"github.com/sirupsen/logrus"
)

func ReportListHandler(tmplSrv service.TemplateService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		templates, err := tmplSrv.ListTemplates()
		if err != nil {
			logrus.Error(err.Error())
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, map[string]interface{}{
			"list": templates,
		})
	}
}
