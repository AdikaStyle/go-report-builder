package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Health() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "UP")
	}
}
