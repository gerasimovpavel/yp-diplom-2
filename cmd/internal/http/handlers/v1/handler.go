package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"yp-diplom-2/cmd/internal/config"
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

func (h *Handler) InitAPI(r chi.Router) {
	ja := jwtauth.New("HS512", []byte(h.cfg.Auth.JWT.SigningKey), nil)
	r.Route("/v1", func(r chi.Router) {
		h.initAuth(r)
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(ja))
			r.Use(jwtauth.Authenticator(ja))
			h.initClubs(r)
		})
	})

}
