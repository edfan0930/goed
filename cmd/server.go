package cmd

import (
	"Cypress/auth/env"
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
		Addr:         ":9453" + env.Port,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	e.Logger.Fatal(e.StartServer(s))
}
