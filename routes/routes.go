package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/muhammadardie/echo-cms/components/abouts"
)

func Register(g *echo.Group) {
	abouts.AboutsRegister(g)
}