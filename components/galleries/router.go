package galleries

import (
	"github.com/labstack/echo/v4"
)

func GalleriesRegister(g *echo.Group) {
	galleries := g.Group("/galleries")
	galleries.GET("", Get)
	galleries.POST("", Create)
	galleries.GET("/:id", Find)
	galleries.PUT("/:id", Update)
	galleries.DELETE("/:id", Destroy)

}
