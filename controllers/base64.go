package controllers

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo"
)

func Base64String(c echo.Context) error {

	fmt.Println("in")
	base64String := c.FormValue("bstring")

	decode, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusOK, err.Error())
	}

	err = ioutil.WriteFile("bb.mp4", decode, 0666)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusOK, err.Error())

	}
	fmt.Println("out")
	return c.JSON(http.StatusOK, decode)
}
