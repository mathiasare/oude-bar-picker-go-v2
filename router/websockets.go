package routes

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"oude-bar-picker-v2/lib"
	"oude-bar-picker-v2/model"
	"oude-bar-picker-v2/service"
	"strconv"

	"github.com/go-chi/chi"
	"gorm.io/gorm"
	"nhooyr.io/websocket"
)

func HandleWebsocketRoutes(ws *lib.WsServer, voteService service.VoteService) *chi.Mux {
	r := chi.NewRouter()
	templates := lib.NewTemplate()

	r.Post("/subscribe/{voteCode}", func(w http.ResponseWriter, r *http.Request) {
		voteCode := chi.URLParam(r, "voteCode")

		// Validate that vote code exists
		if voteCode == "" {
			http.Error(w, "Vote code missing!", http.StatusBadRequest)
			return
		}

		err := ws.Subscribe(r.Context(), w, r, voteCode)
		if errors.Is(err, context.Canceled) {
			return
		}
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
			websocket.CloseStatus(err) == websocket.StatusGoingAway {
			return
		}
		if err != nil {
			log.Println("Error subscribing: ", err)
			return
		}
	})

	r.Get("/subscribe/{voteCode}", func(w http.ResponseWriter, r *http.Request) {
		voteCode := chi.URLParam(r, "voteCode")

		// Validate that vote code exists
		if voteCode == "" {
			http.Error(w, "Vote code missing!", http.StatusBadRequest)
			return
		}

		err := ws.Subscribe(r.Context(), w, r, voteCode)
		if errors.Is(err, context.Canceled) {
			return
		}
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
			websocket.CloseStatus(err) == websocket.StatusGoingAway {
			return
		}
		if err != nil {
			log.Println("Error subscribing: ", err)
			return
		}
	})

	r.Post("/publish/{voteCode}", func(w http.ResponseWriter, r *http.Request) {
		voteCode := chi.URLParam(r, "voteCode")

		// Validate that vote code exists
		if voteCode == "" {
			http.Error(w, "Vote code missing!", http.StatusBadRequest)
			return
		}

		var vDTO model.VoteDTO

		dtoBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("Incorrect form data!")
			http.Error(w, "Incorrect data format!", http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(dtoBytes, &vDTO)
		if err != nil {
			log.Println("Failed to parse json!: ", err)
			http.Error(w, "Failed to parse JSON!", http.StatusBadRequest)
			return
		}

		pId, err := strconv.Atoi(vDTO.ParticipantId)
		if err != nil {
			log.Println("Incorrect form data!")
			http.Error(w, "Incorrect data format!", http.StatusBadRequest)
			return
		}

		barId, err := strconv.Atoi(vDTO.BarId)
		if err != nil {
			log.Println("Incorrect form data!")
			http.Error(w, "Incorrect data format!", http.StatusBadRequest)
			return
		}

		statsData, err := voteService.VoteForBar(uint(pId), uint(barId))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, "User not registered to vote!", http.StatusForbidden)
				return
			}
			http.Error(w, "Failed to update participants vote!", http.StatusInternalServerError)
			return
		}

		temp := templates.Templates.Lookup("stats-row")

		var buf bytes.Buffer
		writer := bufio.NewWriter(&buf)
		temp.Execute(writer, statsData)

		err = writer.Flush()
		if err != nil {
			http.Error(w, "Failed to write websocket response!", http.StatusInternalServerError)
			return
		}

		ws.Publish(voteCode, buf.Bytes())
		w.WriteHeader(http.StatusAccepted)
	})

	r.Post("/subscribed/{voteCode}", func(w http.ResponseWriter, r *http.Request) {
		voteCode := chi.URLParam(r, "voteCode")

		// Validate that vote code exists
		if voteCode == "" {
			http.Error(w, "Vote code missing!", http.StatusBadRequest)
			return
		}

		votePs, err := voteService.PService.GetAllParticipantsForVote(voteCode)
		if err != nil {
			http.Error(w, "Failed to write websocket response!", http.StatusInternalServerError)
			return
		}

		temp := templates.Templates.Lookup("p-row")

		var buf bytes.Buffer
		writer := bufio.NewWriter(&buf)
		temp.Execute(writer, votePs)

		err = writer.Flush()
		if err != nil {
			http.Error(w, "Failed to write websocket response!", http.StatusInternalServerError)
			return
		}

		ws.Publish(voteCode, buf.Bytes())
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return r
}
