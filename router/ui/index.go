package routes

import (
	"net/http"

	"oude-bar-picker-v2/lib"

	"github.com/go-chi/chi"
)

func HandleIndexPage() *chi.Mux {
	r := chi.NewRouter()
	templates := lib.NewTemplate()

	// Index page
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		templates.Templates.ExecuteTemplate(w, "index", nil)
	})

	// Used for static templates with no injected data
	r.Get("/partial/{target}", func(w http.ResponseWriter, r *http.Request) {
		target := chi.URLParam(r, "target")
		templates.Templates.ExecuteTemplate(w, target, nil)
	})

	return r
}
