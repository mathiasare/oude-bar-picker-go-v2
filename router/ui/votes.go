package routes

import (
	"errors"
	"fmt"
	"net/http"
	"oude-bar-picker-v2/lib"
	"oude-bar-picker-v2/model"
	"oude-bar-picker-v2/service"

	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

type VotePageLoadData struct {
	VoteCode string
	Name     string
}

type VotePageData struct {
	Bars         []model.Bar
	VoteCode     string
	Participant  model.Participant
	VotesData    model.VoteStatsDTO
	Participants []model.Participant
}

type EndPageData struct {
	WinningBar   model.Bar
	VoteCode     string
	WinningScore uint
	TotalVotes   int
}

type EndPageEmptyData struct {
	VoteCode string
}

func HandleVotesPage(barService service.BarService, voteService service.VoteService, ws *lib.WsServer) *chi.Mux {
	r := chi.NewRouter()
	templates := lib.NewTemplate()

	// Helper functions
	handleVoteEnded := func(voteCode string, w http.ResponseWriter) {
		vote, err := voteService.GetDeletedById(voteCode)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, "Vote not found!", http.StatusNotFound)
				return
			}
		}

		stats, err := voteService.GetVoteStats(voteCode)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Result page: problem showing the result page.", http.StatusInternalServerError)
			return
		}

		if vote.WinnerId == nil {

			pageData := &EndPageEmptyData{
				VoteCode: voteCode,
			}

			w.WriteHeader(http.StatusOK)
			templates.Templates.ExecuteTemplate(w, "vote-end-empty", pageData)
			return
		}

		winningBar, err := barService.GetById(*vote.WinnerId)
		if err != nil {
			http.Error(w, "End vote: Error fetching bar!", http.StatusInternalServerError)
			return
		}

		var winnerStats model.VoteStatsRow
		for _, s := range stats {
			if s.BarId == *vote.WinnerId {
				winnerStats = s
			}
		}

		totalVotes := 0
		for _, s := range stats {
			totalVotes += int(s.VoteCount)
		}

		pageData := &EndPageData{
			WinningBar:   winningBar,
			VoteCode:     voteCode,
			WinningScore: winnerStats.VoteCount,
			TotalVotes:   totalVotes,
		}

		w.WriteHeader(http.StatusOK)
		templates.Templates.ExecuteTemplate(w, "vote-end", pageData)
	}

	// Endpoints
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		voteCode := r.URL.Query().Get("voteCode")

		// Validate that name and code exists
		if name == "" {
			http.Error(w, "Name missing!", http.StatusBadRequest)
			return
		}

		if voteCode == "" {
			http.Error(w, "Vote code missing!", http.StatusBadRequest)
			return
		}

		loadData := &VotePageLoadData{
			Name:     name,
			VoteCode: voteCode,
		}

		w.WriteHeader(http.StatusOK)
		templates.Templates.ExecuteTemplate(w, "vote", loadData)
	})

	r.Post("/content", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		voteCode := r.URL.Query().Get("voteCode")

		// Validate that name and code exists
		if name == "" {
			http.Error(w, "Name missing!", http.StatusBadRequest)
			return
		}

		if voteCode == "" {
			http.Error(w, "Vote code missing!", http.StatusBadRequest)
			return
		}

		// Make sure that the vote has not ended
		vote, err := voteService.GetById(voteCode)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				handleVoteEnded(voteCode, w)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		votePs := vote.Participants
		var participant *model.Participant
		participant = nil
		for _, p := range votePs {
			if p.Name == name {
				participant = &p
			}
		}

		if participant == nil {
			http.Error(w, "User not registered to vote!", http.StatusForbidden)
			return
		}

		// Fetch bar data
		bars, err := barService.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		votesData, err := voteService.GetVoteStats(voteCode)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pageData := VotePageData{
			Bars:         bars,
			VoteCode:     voteCode,
			Participant:  *participant,
			VotesData:    votesData,
			Participants: votePs,
		}

		w.WriteHeader(http.StatusOK)
		templates.Templates.ExecuteTemplate(w, "vote-page", pageData)
	})

	r.Post("/create", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		name := r.Form.Get("name")

		// Validate that name exists
		if name == "" {
			http.Error(w, "Name missing!", http.StatusBadRequest)
			return
		}

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
		if name == "" {
			http.Error(w, "Name missing!", http.StatusBadRequest)
			return
		}

		if voteCode == "" {
			http.Error(w, "Vote code missing!", http.StatusBadRequest)
			return
		}

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
		r.ParseForm()
		name := r.Form.Get("name")
		voteCode := r.Form.Get("voteCode")

		// Validate that name and code exists
		if name == "" {
			http.Error(w, "Name missing!", http.StatusBadRequest)
			return
		}

		if voteCode == "" {
			http.Error(w, "Vote code missing!", http.StatusBadRequest)
			return
		}

		_, err := voteService.EndVote(voteCode, name)
		if err != nil {
			http.Error(w, "End vote: Error ending vote!", http.StatusInternalServerError)
			return
		}

		ws.Publish(voteCode, []byte("finish"))
	})
	return r
}
