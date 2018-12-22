package router

import (
	"Cypress/gameboy/module/errcode"
	"net/http"
	"strings"

	"github.com/edfan0930/goed/controllers"
	"github.com/labstack/echo"
)

//Router --
func Router(e *echo.Echo) {
	u := e.Group("/user", Auth)

	u.GET("/:uid", controllers.Get)

	u.POST("/:uid", controllers.Create)
}

//Auth --
//middleware
func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		//取token 並判斷是否為空
		token := strings.TrimSpace(c.Request().Header.Get("token"))
		if len(token) == 0 {
			output := errcode.CQ9.Output("4")
			return c.JSON(http.StatusOK, output)
		}
		return next(c)
	}
}
