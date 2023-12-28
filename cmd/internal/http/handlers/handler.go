package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"yp-diplom-2/cmd/internal/config"
	v1 "yp-diplom-2/cmd/internal/http/handlers/v1"
	"yp-diplom-2/cmd/internal/http/response"
	"yp-diplom-2/cmd/internal/service"
)

type Handler struct {
	cfg      *config.Config
	services *service.Services
}

func NewHandler(cfg *config.Config, services *service.Services) *Handler {
	return &Handler{
		cfg:      cfg,
		services: services,
	}
}

func (h *Handler) Init() chi.Router {
	var err error
	r := chi.NewRouter()
	r.Use(
		middleware.Logger,
		SetContentType,
		middleware.AllowContentEncoding("gzip", "deflate"),
		middleware.Compress(5),
		middleware.AllowContentType("application/json"),
		middleware.Recoverer,
	)
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("pong"))
		if err != nil {
			log.Fatal(err)
		}
	})
	h.InitAPI(r)
	return r
}

func (h *Handler) InitAPI(r chi.Router) {
	apiV1 := v1.NewHandler(h.cfg, h.services)
	r.Route("/api", func(r chi.Router) {
		apiV1.InitAPI(r)
	})
	r.MethodNotAllowed(MethodNotAllowed)
	r.NotFound(MethodNotFound)
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	response.NewResponse(w, http.StatusMethodNotAllowed, "method not allowed")
}

func MethodNotFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "method not found", http.StatusNotFound)
}
