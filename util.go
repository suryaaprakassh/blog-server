package main

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func Render(c echo.Context, component templ.Component, statusCode int) error {
	c.Response().Writer.WriteHeader(statusCode)

	err := component.Render(context.Background(), c.Response().Writer)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to render response to template")
	}

	return nil
}
