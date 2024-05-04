package routes

import (
	"encoding/json"
	"errors"
	"log"
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

		barsJson, err := json.Marshal(bars)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(barsJson)
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
				http.Error(w, "Bar not found!", http.StatusNotFound)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		barJson, err := json.Marshal(bar)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(barJson)
	})

	// Create bar
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		res, err := service.Create(r.Body)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(strconv.Itoa(int(res.ID))))
	})

	// Delete bar
	r.Delete("/{barId}", func(w http.ResponseWriter, r *http.Request) {
		barId := chi.URLParam(r, "barId")
		resultId, err := service.Delete(barId)

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, "Bar not found!", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(resultId))
	})

	return r
}
