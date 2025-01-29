package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/miku/cali/internal/config"
	"github.com/miku/cali/internal/db"
	"github.com/miku/cali/internal/models"
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

// Request and response structures
type createAppointmentRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
}

func (s *Server) handleListAppointments(w http.ResponseWriter, r *http.Request) {
	// For now, just return empty list
	s.respondJSON(w, http.StatusOK, []models.Appointment{})
}

func (s *Server) handleCreateAppointment(w http.ResponseWriter, r *http.Request) {
	var req createAppointmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// For now, hardcode user_id as 1
	appt := &models.Appointment{
		UserID:      1,
		Title:       req.Title,
		Description: req.Description,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
	}

	if err := s.db.CreateAppointment(appt); err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to create appointment")
		return
	}

	s.respondJSON(w, http.StatusCreated, appt)
}

func (s *Server) handleGetAppointment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid appointment ID")
		return
	}

	appt, err := s.db.GetAppointment(id)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to get appointment")
		return
	}
	if appt == nil {
		s.respondError(w, http.StatusNotFound, "Appointment not found")
		return
	}

	s.respondJSON(w, http.StatusOK, appt)
}

func (s *Server) handleUpdateAppointment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid appointment ID")
		return
	}

	var req createAppointmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	appt := &models.Appointment{
		ID:          id,
		UserID:      1, // Hardcoded for now
		Title:       req.Title,
		Description: req.Description,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
	}

	if err := s.db.UpdateAppointment(appt); err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to update appointment")
		return
	}

	s.respondJSON(w, http.StatusOK, appt)
}

func (s *Server) handleDeleteAppointment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid appointment ID")
		return
	}

	if err := s.db.DeleteAppointment(id, 1); err != nil { // Hardcoded user_id
		s.respondError(w, http.StatusInternalServerError, "Failed to delete appointment")
		return
	}

	s.respondJSON(w, http.StatusNoContent, nil)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	// For now, just return a simple message
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Welcome to Cali Appointment Scheduler"))
}
