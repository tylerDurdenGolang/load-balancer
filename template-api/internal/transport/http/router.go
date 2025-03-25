package http

import (
	"template-api/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.IService
}

func NewHandler(services service.IService) *Handler {
	return &Handler{service: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger()) // Adding Logger middleware

	api := router.Group("/api")
	{
		items := api.Group("items")
		{
			items.POST("/createItem", h.createItem)
			items.POST("/updateItem", h.updateItem)
			items.GET("/", h.getAllItems)
			items.GET("/:id", h.getItemById)
			items.DELETE("/deleteItem", h.deleteItem)
		}
	}

	return router
}
