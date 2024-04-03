package routes

import (
	"encoding/json"
	"errors"
	"net/http"
	"oude-bar-picker-v2/service"
	"strconv"

	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

func HandleApiBarRoutes(service service.BarService) *chi.Mux {
	r := chi.NewRouter()

	// Get all bars
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		bars, err := service.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonFormat, err := json.Marshal(bars)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(200)
		w.Write(jsonFormat)
	})

	// Get bar by ID
	r.Get("/{barId}", func(w http.ResponseWriter, r *http.Request) {
		barIdStr := chi.URLParam(r, "barId")
		barId, err := strconv.Atoi(barIdStr)
		if err != nil {
			http.Error(w, "Invalid bar id!", http.StatusBadRequest)
			return
		}

		bar, err := service.GetById(uint(barId))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				w.WriteHeader(404)
				w.Write([]byte("Bar not found!"))
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonFormat, err := json.Marshal(bar)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(200)
		w.Write(jsonFormat)
	})

	// Create bar
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		res, err := service.Create(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(201)
		w.Write([]byte(strconv.Itoa(int(res.ID))))
	})

	// Update bar
	r.Put("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("TODO!"))
	})

	// Delete bar
	r.Delete("/{barId}", func(w http.ResponseWriter, r *http.Request) {
		barId := chi.URLParam(r, "barId")
		resultId, err := service.Delete(barId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(200)
		w.Write([]byte(resultId))
	})

	return r
}
