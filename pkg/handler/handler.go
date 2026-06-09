package handler

import (
	"esaccount/pkg/repository"
	"esaccount/pkg/service"
	"net/http"

	"esaccount/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
	repo     *repository.Repository
}

func NewHandler(services *service.Service, repo *repository.Repository) *Handler {
	return &Handler{services: services, repo: repo}
}

//	@title			EasyStartup Account
//	@version		1.0
//	@description	Account service for EasyStartup.
//
//	@contact.name	API Support
//	@contact.email	red.tea.dev@gmail.com
//
//	@host		account.easystartup.su
//	@BasePath	/api
//

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	docs.SwaggerInfo.BasePath = "/"

	router.StaticFS("/static", http.Dir("./public/usercontent"))

	router.GET("/api/public/:id", h.getPublic)
	account := router.Group("/api/account", h.userIdentity)
	{
		account.GET("/", h.getAccount)
		account.POST("/", h.editAccount) //edit
		account.POST("/picture", h.picture)

		//admin
		account.GET("/:id", h.getAccountById)
		account.POST("/:id", h.editAccountById) //edit
		account.DELETE("/:id", h.delete)
	}

	conf := router.Group("/api/confidential", h.userIdentity)
	{
		conf.GET("/", h.getConfidential)
		conf.POST("/", h.setConfidential)
		conf.GET("/:id", h.getConfidentialById)
		conf.POST("/:id", h.setConfidentialById)
	}

	following := router.Group("/api/following", h.userIdentity)
	{
		following.PUT("/follow/:id", h.follow)
		following.PUT("/unfollow/:id", h.unfollow)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}
