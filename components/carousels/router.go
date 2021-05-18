package carousels

import (
	"github.com/labstack/echo/v4"
)

func CarouselsRegister(g *echo.Group) {
	carousels := g.Group("/carousels")
	carousels.GET("", Get)
	carousels.POST("", Create)
	carousels.GET("/:id", Find)
	carousels.PUT("/:id", Update)
	carousels.DELETE("/:id", Destroy)
}
