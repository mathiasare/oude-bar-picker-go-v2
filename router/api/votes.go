package routes

import (
	"github.com/go-chi/chi"
)

func HandleApiVoteRoutes() *chi.Mux {
	r := chi.NewRouter()
	return r
}
