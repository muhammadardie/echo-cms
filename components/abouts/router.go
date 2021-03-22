package abouts

import (
	"github.com/labstack/echo/v4"
)

func AboutsRegister(g *echo.Group) {
	abouts := g.Group("/abouts")
	abouts.GET("", Get)
	abouts.POST("", Create)
}