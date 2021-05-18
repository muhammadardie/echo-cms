package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/muhammadardie/echo-cms/components/abouts"
	"github.com/muhammadardie/echo-cms/components/blogs"
	"github.com/muhammadardie/echo-cms/components/carousels"
	"github.com/muhammadardie/echo-cms/components/companies"
	"github.com/muhammadardie/echo-cms/components/contacts"
	"github.com/muhammadardie/echo-cms/components/galleries"
	"github.com/muhammadardie/echo-cms/components/headers"
	"github.com/muhammadardie/echo-cms/components/services"
	"github.com/muhammadardie/echo-cms/components/socmeds"
	"github.com/muhammadardie/echo-cms/components/teams"
	"github.com/muhammadardie/echo-cms/components/testimonies"
	"github.com/muhammadardie/echo-cms/components/users"
	"net/http"
)

func Register(g *echo.Group) {
	g.GET("/csrf", func(c echo.Context) error {
		return c.JSON(http.StatusOK, c.Get("csrf"))
	})
	abouts.AboutsRegister(g)
	blogs.BlogsRegister(g)
	carousels.CarouselsRegister(g)
	companies.CompaniesRegister(g)
	contacts.ContactsRegister(g)
	galleries.GalleriesRegister(g)
	headers.HeadersRegister(g)
	services.ServicesRegister(g)
	socmeds.SocmedsRegister(g)
	teams.TeamsRegister(g)
	testimonies.TestimoniesRegister(g)
	users.UsersRegister(g)
}
