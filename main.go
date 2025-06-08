package main

import (
	"context"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/suryaaprakassh/blog_server/posts"
	"github.com/suryaaprakassh/blog_server/templates"
)

var blogs map[string]bool

func GenerateStatic(p *posts.Post, path string) error {
	path = strings.TrimSuffix(path, "/")
	file, err := os.Create(path + "/" + p.Slug + ".html")
	if err != nil {
		return err
	}
	defer file.Close()
	err = templates.Index(p.Title, p.Date, p.Content).Render(context.Background(), file)
	if err != nil {
		return err
	}
	return file.Sync()
}

func generateHome(post []*posts.Post) {
	file, err := os.Create("./static/home.html")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	err = templates.Home(post).Render(context.Background(), file)
	if err != nil {
		panic(err)
	}
	if err := file.Sync(); err != nil {
		panic(err)
	}
}

func init() {
	dir := "./posts"
	blogs = make(map[string]bool)
	root := os.DirFS(dir)
	mdFiles, err := fs.Glob(root, "*.md")
	postsContainer := make([]*posts.Post, 0)
	if err != nil {
		panic(err)
	}
	for _, path := range mdFiles {
		f, err := os.Open(dir + "/" + path)
		if err != nil {
			panic(err)
		}
		p := posts.NewPost(f)
		err = GenerateStatic(p, "./static")
		if err != nil {
			panic(err)
		}
		postsContainer = append(postsContainer, p)
		blogs[p.Slug] = true
	}
	generateHome(postsContainer)
}

func keepAlive(t *time.Ticker) {
	for {
		select {
			case <-t.C: 
			res,err := http.Get("http://localhost:8000/health")
			if err != nil {
					slog.Error(err.Error())
			}
			slog.Info("[PING]", "code",res.StatusCode)
		}
	}
}

func main() {
	e := echo.New()
	t := time.NewTicker(time.Minute * 10);
	go keepAlive(t);
	e.Static("/static", "./static")
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK,"ok")
	})
	e.GET("/",func(c echo.Context) error {
		return c.File("static/home.html")
	})
	e.GET("/:path", func(c echo.Context) error {
		slug := c.Param("path")
		if _, ok := blogs[slug]; !ok {
			return Render(c, templates.Notfound(), http.StatusNotFound)
		}
		return c.File("static/" + slug + ".html")
	})

	e.Logger.Fatal(e.Start(":8000"))
}
