package companies

import (
	"github.com/labstack/echo/v4"
)

func CompaniesRegister(g *echo.Group) {
	companies := g.Group("/companies")
	companies.GET("", Get)
	companies.POST("", Create)
	companies.GET("/:id", Find)
	companies.PUT("/:id", Update)
	companies.DELETE("/:id", Destroy)
}
