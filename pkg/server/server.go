package server

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"tigerhall-kittens-app/pkg/auth"
	"tigerhall-kittens-app/pkg/handlers"
	"tigerhall-kittens-app/pkg/middleware"
	"tigerhall-kittens-app/pkg/service"
)

type server struct {
	router *mux.Router
	logger *log.Logger
}

func NewServer() *server {
	return &server{
		router: mux.NewRouter(),
		logger: log.New(os.Stdout, "[Tigerhall Kittens] ", log.LstdFlags),
	}
}

func (s *server) SetupRoutes(tigerService service.TigerService, auth *auth.Auth) {
	handlers := handlers.NewHandlers(tigerService, s.logger, auth)

	// Public routes
	s.router.HandleFunc("/signup", handlers.SignupHandler).Methods("POST")
	s.router.HandleFunc("/login", handlers.LoginHandler).Methods("POST")

	s.router.HandleFunc("/tigers", handlers.GetAllTigersHandler).Methods("GET")
	s.router.HandleFunc("/tiger/{id}/sightings", handlers.GetTigerSightingsByIDHandler).Methods("GET")

	// Protected routes (require authentication)
	s.router.Handle("/tiger/create", middleware.AuthMiddleware(auth, http.HandlerFunc(handlers.CreateTigerHandler))).Methods("POST")
	s.router.Handle("/tiger-sighting/create", middleware.AuthMiddleware(auth, http.HandlerFunc(handlers.CreateTigerSightingHandler))).Methods("POST")
}

func (s *server) Start(port string) error {
	s.logger.Printf("Starting server on port %s...", port)
	return http.ListenAndServe(":"+port, s.router)
}
