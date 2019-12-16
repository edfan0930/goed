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
		Data   []count `json:"data"`
		Status status  `json:"status"`
	}
	login struct {
		Data   user   `json:"data"`
		Status status `json:"status"`
	}

	count struct {
		Name      string `json:"name"`
		ID        int64  `json:"id"`
		StartNo   int64  `json:"start_no"`
		EndNo     int64  `json:"end_no"`
		Wait      int64  `json:"wait"`
		CurrentNo int64  `json:"current_no"`
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
	a := count{"aaa", 1, 1, 99, 1, 5}
	b := count{"bbb", 2, 100, 199, 5, 111}
	su := []count{a, b}
	return c.JSON(http.StatusOK, response{Data: su})
}

//Login --
func Login(c echo.Context) error {
	acc := c.FormValue("account")
	pas := c.FormValue("password")
	if acc == "" || pas == "" {

		return c.JSON(http.StatusOK, response{Status: status{101, "Not Found"}})
	}
	return c.JSON(http.StatusOK, login{Data: user{"ed", "e", "male"}, Status: status{0, "Success"}})
}
