package handlers

import (
	"github.com/AdikaStyle/go-report-builder/internal/business"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func ReportListHandler(tmplSrv business.TemplateService) gin.HandlerFunc {
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
