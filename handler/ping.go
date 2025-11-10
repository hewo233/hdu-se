package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hewo233/hdu-se/models"
	"net/http"
)

type PingResponse struct {
	Message string `json:"message"`
}

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, models.Report{
		Code: http.StatusOK,
		Result: PingResponse{
			Message: "pong",
		},
	})
}
