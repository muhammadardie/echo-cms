package teams

import (
	"github.com/labstack/echo/v4"
)

func TeamsRegister(g *echo.Group) {
	teams := g.Group("/teams")
	teams.GET("", Get)
	teams.POST("", Create)
	teams.GET("/:id", Find)
	teams.PUT("/:id", Update)
	teams.DELETE("/:id", Destroy)
}
