package router

import (
	"github.com/edfan0930/goed/controllers"
	"github.com/labstack/echo"
)

//Router --
func Router(e *echo.Echo) {

	e.POST("login", controllers.Login)

	u := e.Group("/count", Auth)

	u.GET("", controllers.Get)

	u.POST("/:uid", controllers.Create)

	e.POST("/bstring", controllers.Base64String)
}

//Auth --
//middleware
func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		//取token 並判斷是否為空
		/* 	token := strings.TrimSpace(c.Request().Header.Get("token"))
		if len(token) == 0 {
			return c.JSON(http.StatusOK, "")
		} */
		return next(c)
	}
}
