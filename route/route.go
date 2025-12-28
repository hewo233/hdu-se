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

	auth := R.Group("/auth")
	auth.POST("/register", handler.RegisterUser)
	auth.POST("/login", handler.UserLogin)

	user := R.Group("/user")
	user.Use(middleware.JWTAuth("user"))
	user.GET("/:id", handler.GetUserInfoByID)
	user.GET("", handler.GetUserInfoByEmail)

	coze := R.Group("/coze")
	coze.Use(middleware.JWTAuth("user"))
	coze.POST("/conversation", handler.CreateConversation)
	coze.GET("/conversation", handler.ListConversations)
	coze.POST("/chat", handler.CreateChat)
	coze.GET("/chat", handler.RetrieveConversation)
	coze.GET("/chat/message", handler.ChatMessageList)
	coze.GET("/conversation/message", handler.ConversationMessageList)
}
