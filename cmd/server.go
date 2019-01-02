package cmd

import (
	"net/http"
	"time"

	"github.com/edfan0930/goed/router"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func server() {

	e := echo.New()

	e.Use(middleware.Recover())

	router.Router(e)

	s := &http.Server{
		Addr:         "127.0.0.1:9453",
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	e.Logger.Fatal(e.StartServer(s))
}
