package lib

import (
	"html/template"
	"io"

	"github.com/go-chi/chi"
)

type Templates struct {
	Templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c chi.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}

func NewTemplate() *Templates {
	return &Templates{
		Templates: template.Must(template.ParseGlob("views/*.html")),
	}
}
