package templates

const BaseServer = `package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Routes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", Index)

	return e
}

func NewServer() *http.Server {
	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", 8000),
		Handler:      Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return s
}

func Render(e echo.Context, code int, c templ.Component) error {
	e.Response().Status = code
	return c.Render(e.Request().Context(), e.Response().Writer)
}
`
