package routes

import (
	"errors"
	"log"
	"net/http"
	"oude-bar-picker-v2/lib"
	"oude-bar-picker-v2/model"
	"oude-bar-picker-v2/service"

	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

type VotePageData struct {
	Bars        []model.Bar
	VoteCode    string
	Participant model.Participant
	TableData   model.VoteStatsDTO
}

type EndPageData struct {
	WinningBar     model.Bar
	VoteCode       string
	FinishedByName string
	WinningScore   uint
	TotalVotes     int
}

type EndPageEmptyData struct {
	VoteCode       string
	FinishedByName string
}

func HandleVotesPage(barService service.BarService, voteService service.VoteService, ws *service.WsServer) *chi.Mux {
	r := chi.NewRouter()
	templates := lib.NewTemplate()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		voteCode := r.URL.Query().Get("voteCode")

		_, err := voteService.GetById(voteCode)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, "Vote not found!", http.StatusNotFound)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		participant, err := voteService.PService.GetByNameAndVote(name, voteCode)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, "User not registered to vote!", http.StatusForbidden)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		bars, err := barService.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		votePs, err := voteService.PService.Repo.FindAllForVoteWhereBarNotNull(voteCode)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tableData := voteService.PService.GetVoteStats(votePs)
		pageData := VotePageData{
			Bars:        bars,
			VoteCode:    voteCode,
			Participant: participant,
			TableData:   tableData,
		}

		log.Println("Subscribed to websocket server!")
		templates.Templates.ExecuteTemplate(w, "vote", pageData)
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		name := r.Form.Get("name")

		// Validate that name exists

		//Create new vote
		vote, err := voteService.Create(model.Vote{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Add user to vote
		err = voteService.AddUserToVote(vote.ID, name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect
		url := "/vote?name=" + name + "&voteCode=" + vote.ID
		http.Redirect(w, r, url, http.StatusMovedPermanently)
	})

	r.Post("/join", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		name := r.Form.Get("name")
		voteCode := r.Form.Get("voteCode")

		// Validate that name and code exists

		// Update vote with new user
		err := voteService.AddUserToVote(voteCode, name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		url := "/vote?name=" + name + "&voteCode=" + voteCode
		http.Redirect(w, r, url, http.StatusMovedPermanently)
	})

	r.Post("/end", func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		voteCode := chi.URLParam(r, "voteCode")

		// Validate that name and code exists

		stats, err := voteService.EndVote(voteCode, name)
		if err != nil {
			http.Error(w, "End vote: Error ending vote!", http.StatusInternalServerError)
			return
		}

		if len(stats) == 0 {
			pageData := &EndPageEmptyData{
				VoteCode:       voteCode,
				FinishedByName: name,
			}
			templates.Templates.ExecuteTemplate(w, "end-vote-empty", pageData)
			return
		}

		winnerStats := stats[0]
		winningBar, err := barService.GetById(winnerStats.BarId)
		if err != nil {
			http.Error(w, "End vote: Error fetching bar!", http.StatusInternalServerError)
			return
		}

		pageData := &EndPageData{
			WinningBar:     winningBar,
			VoteCode:       voteCode,
			FinishedByName: name,
			WinningScore:   winnerStats.VoteCount,
			TotalVotes:     len(stats),
		}

		templates.Templates.ExecuteTemplate(w, "end-vote", pageData)
	})

	return r
}
