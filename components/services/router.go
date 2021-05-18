package services

import (
	"github.com/labstack/echo/v4"
)

func ServicesRegister(g *echo.Group) {
	services := g.Group("/services")
	services.GET("", Get)
	services.POST("", Create)
	services.GET("/:id", Find)
	services.PUT("/:id", Update)
	services.DELETE("/:id", Destroy)

}
