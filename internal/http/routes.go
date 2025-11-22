package http

import (
	"log"
	"net/http"

	"antonvedaet/internship_task/internal/http/handlers"
	"antonvedaet/internship_task/internal/service"
	"antonvedaet/internship_task/internal/store"
)

func MakeMux() *http.ServeMux {
	mux := http.NewServeMux()

	db, err := store.New()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	teamService := service.NewTeamService(db)
	userService := service.NewUserService(db)
	prService := service.NewPRService(db)

	handler := handlers.NewHandlers(
		teamService,
		userService,
		prService,
	)
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
	mux.HandleFunc("POST /team/add", handler.AddTeam)
	mux.HandleFunc("GET /team/get", handler.GetTeam)
	mux.HandleFunc("POST /team/deactivate", handler.DeactivateTeamUsers)

	mux.HandleFunc("POST /users/setIsActive", handler.SetUserActive)
	mux.HandleFunc("GET /users/getReview", handler.GetUserReview)

	mux.HandleFunc("POST /pullRequest/create", handler.CreatePR)
	mux.HandleFunc("POST /pullRequest/merge", handler.MergePR)
	mux.HandleFunc("POST /pullRequest/reassign", handler.ReassignReviewer)

	mux.HandleFunc("GET /health", handler.Health)

	return mux
}
