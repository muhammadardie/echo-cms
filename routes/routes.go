package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/muhammadardie/echo-cms/components/abouts"
	"github.com/muhammadardie/echo-cms/components/blogs"
	"net/http"
)

func Register(g *echo.Group) {
	g.GET("/csrf", func(c echo.Context) error {
		return c.JSON(http.StatusOK, c.Get("csrf"))
	})
	abouts.AboutsRegister(g)
	blogs.BlogsRegister(g)
}
