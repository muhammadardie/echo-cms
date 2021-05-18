package socmeds

import (
	"github.com/labstack/echo/v4"
)

func SocmedsRegister(g *echo.Group) {
	socmeds := g.Group("/socmeds")
	socmeds.GET("", Get)
	socmeds.POST("", Create)
	socmeds.GET("/:id", Find)
	socmeds.PUT("/:id", Update)
	socmeds.DELETE("/:id", Destroy)

}
