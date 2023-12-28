package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"net/http"
)

func (h *Handler) initClubs(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(cors.Handler(cors.Options{
			AllowedMethods: []string{"GET"},
			MaxAge:         300,
		}))
		r.Route("/clubs", func(r chi.Router) {
			r.Get("/", h.loadClubs)
		})
	})
	r.Group(func(r chi.Router) {
		r.Use(cors.Handler(cors.Options{
			AllowedMethods: []string{"POST"},
			MaxAge:         300,
		}))
		r.Route("/club", func(r chi.Router) {
			r.Post("/", h.createClub)
		})
	})
	r.Group(func(r chi.Router) {
		r.Use(cors.Handler(cors.Options{
			AllowedMethods: []string{"GET", "PATCH", "DELETE"},
			MaxAge:         300,
		}))
		r.Route("/club/{clubId}", func(r chi.Router) {
			r.Get("/", h.loadClub)
		})
	})
}

func (h *Handler) loadClubs(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) createClub(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) loadClub(w http.ResponseWriter, r *http.Request) {

}
