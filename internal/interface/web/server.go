package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joshthewhite/poolvibes/internal/application/services"
	"github.com/joshthewhite/poolvibes/internal/interface/web/handlers"
)

type Server struct {
	mux       *http.ServeMux
	authSvc   *services.AuthService
	userSvc   *services.UserService
	chemSvc   *services.ChemistryService
	taskSvc   *services.TaskService
	equipSvc  *services.EquipmentService
	chemicSvc *services.ChemicalService
}

func NewServer(authSvc *services.AuthService, userSvc *services.UserService, chemSvc *services.ChemistryService, taskSvc *services.TaskService, equipSvc *services.EquipmentService, chemicSvc *services.ChemicalService) *Server {
	s := &Server{
		mux:       http.NewServeMux(),
		authSvc:   authSvc,
		userSvc:   userSvc,
		chemSvc:   chemSvc,
		taskSvc:   taskSvc,
		equipSvc:  equipSvc,
		chemicSvc: chemicSvc,
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	pageHandler := handlers.NewPageHandler()
	authHandler := handlers.NewAuthHandler(s.authSvc)
	chemHandler := handlers.NewChemistryHandler(s.chemSvc)
	taskHandler := handlers.NewTaskHandler(s.taskSvc)
	equipHandler := handlers.NewEquipmentHandler(s.equipSvc)
	chemicHandler := handlers.NewChemicalHandler(s.chemicSvc)
	adminHandler := handlers.NewAdminHandler(s.userSvc)
	settingsHandler := handlers.NewSettingsHandler(s.userSvc)

	auth := func(h http.HandlerFunc) http.HandlerFunc { return requireAuth(s.authSvc, h) }
	admin := func(h http.HandlerFunc) http.HandlerFunc { return requireAdmin(s.authSvc, h) }

	// Auth routes (no auth required)
	s.mux.HandleFunc("GET /login", authHandler.LoginPage)
	s.mux.HandleFunc("POST /login", authHandler.Login)
	s.mux.HandleFunc("GET /signup", authHandler.SignupPage)
	s.mux.HandleFunc("POST /signup", authHandler.Signup)
	s.mux.HandleFunc("POST /logout", authHandler.Logout)

	// Page (auth required)
	s.mux.HandleFunc("GET /{$}", auth(pageHandler.Index))

	// Chemistry (auth required)
	s.mux.HandleFunc("GET /chemistry", auth(chemHandler.List))
	s.mux.HandleFunc("GET /chemistry/new", auth(chemHandler.NewForm))
	s.mux.HandleFunc("POST /chemistry", auth(chemHandler.Create))
	s.mux.HandleFunc("GET /chemistry/{id}/edit", auth(chemHandler.EditForm))
	s.mux.HandleFunc("PUT /chemistry/{id}", auth(chemHandler.Update))
	s.mux.HandleFunc("DELETE /chemistry/{id}", auth(chemHandler.Delete))

	// Tasks (auth required)
	s.mux.HandleFunc("GET /tasks", auth(taskHandler.List))
	s.mux.HandleFunc("GET /tasks/new", auth(taskHandler.NewForm))
	s.mux.HandleFunc("POST /tasks", auth(taskHandler.Create))
	s.mux.HandleFunc("GET /tasks/{id}/edit", auth(taskHandler.EditForm))
	s.mux.HandleFunc("PUT /tasks/{id}", auth(taskHandler.Update))
	s.mux.HandleFunc("POST /tasks/{id}/complete", auth(taskHandler.Complete))
	s.mux.HandleFunc("DELETE /tasks/{id}", auth(taskHandler.Delete))

	// Equipment (auth required)
	s.mux.HandleFunc("GET /equipment", auth(equipHandler.List))
	s.mux.HandleFunc("GET /equipment/new", auth(equipHandler.NewForm))
	s.mux.HandleFunc("POST /equipment", auth(equipHandler.Create))
	s.mux.HandleFunc("GET /equipment/{id}/edit", auth(equipHandler.EditForm))
	s.mux.HandleFunc("PUT /equipment/{id}", auth(equipHandler.Update))
	s.mux.HandleFunc("DELETE /equipment/{id}", auth(equipHandler.Delete))
	s.mux.HandleFunc("GET /equipment/{id}/service-records/new", auth(equipHandler.NewServiceRecordForm))
	s.mux.HandleFunc("POST /equipment/{id}/service-records", auth(equipHandler.CreateServiceRecord))
	s.mux.HandleFunc("DELETE /equipment/{id}/service-records/{recordId}", auth(equipHandler.DeleteServiceRecord))

	// Chemicals (auth required)
	s.mux.HandleFunc("GET /chemicals", auth(chemicHandler.List))
	s.mux.HandleFunc("GET /chemicals/new", auth(chemicHandler.NewForm))
	s.mux.HandleFunc("POST /chemicals", auth(chemicHandler.Create))
	s.mux.HandleFunc("GET /chemicals/{id}/edit", auth(chemicHandler.EditForm))
	s.mux.HandleFunc("PUT /chemicals/{id}", auth(chemicHandler.Update))
	s.mux.HandleFunc("POST /chemicals/{id}/adjust", auth(chemicHandler.AdjustStock))
	s.mux.HandleFunc("DELETE /chemicals/{id}", auth(chemicHandler.Delete))

	// Settings (auth required)
	s.mux.HandleFunc("GET /settings", auth(settingsHandler.Page))
	s.mux.HandleFunc("PUT /settings", auth(settingsHandler.Update))

	// Admin (admin required)
	s.mux.HandleFunc("GET /admin/users", admin(adminHandler.ListUsers))
	s.mux.HandleFunc("GET /admin/users/{id}/edit", admin(adminHandler.EditUser))
	s.mux.HandleFunc("PUT /admin/users/{id}", admin(adminHandler.UpdateUser))
}

func (s *Server) Start(addr string) error {
	fmt.Printf("PoolVibes server starting on %s\n", addr)
	handler := logRequests(s.mux)
	return http.ListenAndServe(addr, handler)
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
