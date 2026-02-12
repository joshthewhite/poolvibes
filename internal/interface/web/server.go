package web

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/josh/poolio/internal/application/services"
	"github.com/josh/poolio/internal/interface/web/handlers"
)

//go:embed templates/layout.html
var layoutHTML embed.FS

type Server struct {
	router    chi.Router
	chemSvc   *services.ChemistryService
	taskSvc   *services.TaskService
	equipSvc  *services.EquipmentService
	chemicSvc *services.ChemicalService
}

func NewServer(chemSvc *services.ChemistryService, taskSvc *services.TaskService, equipSvc *services.EquipmentService, chemicSvc *services.ChemicalService) *Server {
	s := &Server{
		router:    chi.NewRouter(),
		chemSvc:   chemSvc,
		taskSvc:   taskSvc,
		equipSvc:  equipSvc,
		chemicSvc: chemicSvc,
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)

	pageHandler := handlers.NewPageHandler(layoutHTML)
	chemHandler := handlers.NewChemistryHandler(s.chemSvc)
	taskHandler := handlers.NewTaskHandler(s.taskSvc)
	equipHandler := handlers.NewEquipmentHandler(s.equipSvc)
	chemicHandler := handlers.NewChemicalHandler(s.chemicSvc)

	s.router.Get("/", pageHandler.Index)

	s.router.Route("/chemistry", func(r chi.Router) {
		r.Get("/", chemHandler.List)
		r.Get("/new", chemHandler.NewForm)
		r.Post("/", chemHandler.Create)
		r.Get("/{id}/edit", chemHandler.EditForm)
		r.Put("/{id}", chemHandler.Update)
		r.Delete("/{id}", chemHandler.Delete)
	})

	s.router.Route("/tasks", func(r chi.Router) {
		r.Get("/", taskHandler.List)
		r.Get("/new", taskHandler.NewForm)
		r.Post("/", taskHandler.Create)
		r.Get("/{id}/edit", taskHandler.EditForm)
		r.Put("/{id}", taskHandler.Update)
		r.Post("/{id}/complete", taskHandler.Complete)
		r.Delete("/{id}", taskHandler.Delete)
	})

	s.router.Route("/equipment", func(r chi.Router) {
		r.Get("/", equipHandler.List)
		r.Get("/new", equipHandler.NewForm)
		r.Post("/", equipHandler.Create)
		r.Get("/{id}/edit", equipHandler.EditForm)
		r.Put("/{id}", equipHandler.Update)
		r.Delete("/{id}", equipHandler.Delete)
		r.Get("/{id}/service-records/new", equipHandler.NewServiceRecordForm)
		r.Post("/{id}/service-records", equipHandler.CreateServiceRecord)
		r.Delete("/{id}/service-records/{recordId}", equipHandler.DeleteServiceRecord)
	})

	s.router.Route("/chemicals", func(r chi.Router) {
		r.Get("/", chemicHandler.List)
		r.Get("/new", chemicHandler.NewForm)
		r.Post("/", chemicHandler.Create)
		r.Get("/{id}/edit", chemicHandler.EditForm)
		r.Put("/{id}", chemicHandler.Update)
		r.Post("/{id}/adjust", chemicHandler.AdjustStock)
		r.Delete("/{id}", chemicHandler.Delete)
	})
}

func (s *Server) Start(addr string) error {
	fmt.Printf("PoolVibes server starting on %s\n", addr)
	return http.ListenAndServe(addr, s.router)
}
