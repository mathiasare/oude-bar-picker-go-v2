package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"

	// Import the routes package
	database "oude-bar-picker-v2/db"
	"oude-bar-picker-v2/lib"
	"oude-bar-picker-v2/repository"
	router "oude-bar-picker-v2/router"
	api "oude-bar-picker-v2/router/api"
	ui "oude-bar-picker-v2/router/ui"
	"oude-bar-picker-v2/service"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading environment variables: %s", err)
	}

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)

	// Websockets
	ws := lib.NewWsServer()

	// Initialize db
	db := database.Connect()

	// Seeding
	//database.Seed(db)

	// Initialize components
	barRepository := repository.BarRepository{Db: db}
	barService := service.BarService{Repo: barRepository}
	participantRepository := repository.ParticipantRepository{Db: db}
	participantService := service.ParticipantService{Repo: participantRepository}
	voteRepository := repository.VoteRepository{Db: db}
	voteService := service.VoteService{Repo: voteRepository, PService: participantService}

	// Setup static file serving
	fileServer := http.FileServer(http.Dir("./resource/static/"))
	r.Handle("/resource/static/*", http.StripPrefix("/resource/static/", fileServer))

	// UI route handlers
	r.Mount("/", ui.HandleIndexPage())
	r.Mount("/vote", ui.HandleVotesPage(barService, voteService, ws))

	// Websocket route handlers
	r.Mount("/ws", router.HandleWebsocketRoutes(ws, voteService))

	// API route handlers
	r.Get("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is the oude bar picker backend API."))
	})
	r.Mount("/api/bars", api.HandleApiBarRoutes(barService))
	r.Mount("/api/votes", api.HandleApiVoteRoutes())

	log.Println("Starting server on :3000")
	http.ListenAndServe(":3000", r)
}
