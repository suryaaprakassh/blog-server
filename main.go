package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/suryaaprakassh/blog_server/templates"
)

func main() {
	e := echo.New()
	e.Static("/static", "./static")
	e.GET("/", func(c echo.Context) error {
		return Render(c,templates.Index("Hello","Jun 10,2025",templates.TestContent()),http.StatusOK);
	})


	e.Logger.Fatal(e.Start(":8000"))
}
