package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/miku/cali/internal/config"
	"github.com/miku/cali/internal/db"
)

type Server struct {
	Router *mux.Router
	db     *db.Database
	config *config.Config
}

func NewServer(db *db.Database, cfg *config.Config) *Server {
	s := &Server{
		Router: mux.NewRouter(),
		db:     db,
		config: cfg,
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	// API routes
	api := s.Router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/appointments", s.handleListAppointments).Methods("GET")
	api.HandleFunc("/appointments", s.handleCreateAppointment).Methods("POST")
	api.HandleFunc("/appointments/{id}", s.handleGetAppointment).Methods("GET")
	api.HandleFunc("/appointments/{id}", s.handleUpdateAppointment).Methods("PUT")
	api.HandleFunc("/appointments/{id}", s.handleDeleteAppointment).Methods("DELETE")

	// Web interface routes
	s.Router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/",
			http.FileServer(http.Dir(s.config.Web.StaticDir))))
	s.Router.HandleFunc("/", s.handleIndex).Methods("GET")
}

func (s *Server) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) respondError(w http.ResponseWriter, status int, message string) {
	s.respondJSON(w, status, map[string]string{"error": message})
}
