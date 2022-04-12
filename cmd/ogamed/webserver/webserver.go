package webserver

import (
	"io"
	"log"
	"net/http"
	"text/template"

	"github.com/faunX/ogame/cmd/ogamed/webserver/bindata"
	"github.com/labstack/echo"
)

func Start() *echo.Echo {
	e := echo.New()

	templates, err := bindata.GetDir("template")
	if err != nil {
		log.Println(err)
	}

	t2 := template.New("")

	for _, file := range templates {
		log.Println("Read " + file.Name())
		data, err := bindata.GetFile("template/" + file.Name())
		if err != nil {
			log.Println(err)
		}
		t2.New(file.Name()).Parse(string(data))
	}

	t := &Template{
		templates: t2,
	}

	e.Renderer = t
	e.GET("/hello", Hello)
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "home", "World")
	})
	return e
}

func Hello(c echo.Context) error {
	return c.Render(http.StatusOK, "hello", "World")
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
