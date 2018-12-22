package controllers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

type (
	user struct {
		Name     string `form:"name"`
		Nickname string `form:"nickname"`
		Gender   string `form:"gender"`
	}
	status struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
	}
	response struct {
		Data   []user `json:"data"`
		Status status `json:"status"`
	}
)

// Create --
func Create(c echo.Context) error {
	f, err := c.FormParams()
	fmt.Println("form", f)
	return c.JSON(http.StatusOK, err)
}

// Get --
func Get(c echo.Context) error {
	a := user{"aaaaa", "a", "male"}
	b := user{"bbbb", "b", "female"}
	su := []user{a, b}
	return c.JSON(http.StatusOK, response{Data: su})
}
