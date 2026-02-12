package web

import (
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/joshthewhite/poolvibes/internal/application/services"
	"github.com/joshthewhite/poolvibes/internal/interface/web/handlers"
)

//go:embed templates/layout.html
var layoutHTML embed.FS

type Server struct {
	mux       *http.ServeMux
	chemSvc   *services.ChemistryService
	taskSvc   *services.TaskService
	equipSvc  *services.EquipmentService
	chemicSvc *services.ChemicalService
}

func NewServer(chemSvc *services.ChemistryService, taskSvc *services.TaskService, equipSvc *services.EquipmentService, chemicSvc *services.ChemicalService) *Server {
	s := &Server{
		mux:       http.NewServeMux(),
		chemSvc:   chemSvc,
		taskSvc:   taskSvc,
		equipSvc:  equipSvc,
		chemicSvc: chemicSvc,
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	pageHandler := handlers.NewPageHandler(layoutHTML)
	chemHandler := handlers.NewChemistryHandler(s.chemSvc)
	taskHandler := handlers.NewTaskHandler(s.taskSvc)
	equipHandler := handlers.NewEquipmentHandler(s.equipSvc)
	chemicHandler := handlers.NewChemicalHandler(s.chemicSvc)

	s.mux.HandleFunc("GET /{$}", pageHandler.Index)

	s.mux.HandleFunc("GET /chemistry", chemHandler.List)
	s.mux.HandleFunc("GET /chemistry/new", chemHandler.NewForm)
	s.mux.HandleFunc("POST /chemistry", chemHandler.Create)
	s.mux.HandleFunc("GET /chemistry/{id}/edit", chemHandler.EditForm)
	s.mux.HandleFunc("PUT /chemistry/{id}", chemHandler.Update)
	s.mux.HandleFunc("DELETE /chemistry/{id}", chemHandler.Delete)

	s.mux.HandleFunc("GET /tasks", taskHandler.List)
	s.mux.HandleFunc("GET /tasks/new", taskHandler.NewForm)
	s.mux.HandleFunc("POST /tasks", taskHandler.Create)
	s.mux.HandleFunc("GET /tasks/{id}/edit", taskHandler.EditForm)
	s.mux.HandleFunc("PUT /tasks/{id}", taskHandler.Update)
	s.mux.HandleFunc("POST /tasks/{id}/complete", taskHandler.Complete)
	s.mux.HandleFunc("DELETE /tasks/{id}", taskHandler.Delete)

	s.mux.HandleFunc("GET /equipment", equipHandler.List)
	s.mux.HandleFunc("GET /equipment/new", equipHandler.NewForm)
	s.mux.HandleFunc("POST /equipment", equipHandler.Create)
	s.mux.HandleFunc("GET /equipment/{id}/edit", equipHandler.EditForm)
	s.mux.HandleFunc("PUT /equipment/{id}", equipHandler.Update)
	s.mux.HandleFunc("DELETE /equipment/{id}", equipHandler.Delete)
	s.mux.HandleFunc("GET /equipment/{id}/service-records/new", equipHandler.NewServiceRecordForm)
	s.mux.HandleFunc("POST /equipment/{id}/service-records", equipHandler.CreateServiceRecord)
	s.mux.HandleFunc("DELETE /equipment/{id}/service-records/{recordId}", equipHandler.DeleteServiceRecord)

	s.mux.HandleFunc("GET /chemicals", chemicHandler.List)
	s.mux.HandleFunc("GET /chemicals/new", chemicHandler.NewForm)
	s.mux.HandleFunc("POST /chemicals", chemicHandler.Create)
	s.mux.HandleFunc("GET /chemicals/{id}/edit", chemicHandler.EditForm)
	s.mux.HandleFunc("PUT /chemicals/{id}", chemicHandler.Update)
	s.mux.HandleFunc("POST /chemicals/{id}/adjust", chemicHandler.AdjustStock)
	s.mux.HandleFunc("DELETE /chemicals/{id}", chemicHandler.Delete)
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
