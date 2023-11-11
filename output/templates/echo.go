package templates

const Echo = `package handlers

import (
	"%s/views"

	"github.com/labstack/echo/v4"
)

func Index(e echo.Context) error {
	return Render(e, 200, views.Index("%s"))
}
`
