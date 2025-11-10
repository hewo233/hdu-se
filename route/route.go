package route

import (
	"github.com/gin-gonic/gin"
	"github.com/hewo233/hdu-se/handler"
	"github.com/hewo233/hdu-se/middleware"
)

var R *gin.Engine

func InitRoute() {
	R = gin.New()
	R.Use(gin.Logger(), gin.Recovery())
	R.Use(middleware.CorsMiddleware())

	R.GET("/ping", handler.Ping)
}
