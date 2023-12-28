package v1

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"yp-diplom-2/cmd/internal/domain"
	"yp-diplom-2/cmd/internal/http/response"
	"yp-diplom-2/cmd/internal/service"
)

func (h *Handler) initAuth(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Post("/signup", h.SignUp)
		r.Post("/signin", h.SignIn)
	})
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var inp service.SignUpInput

	if err := json.NewDecoder(r.Body).Decode(&inp); err != nil {
		response.NewResponse(w, http.StatusBadRequest, "invalid input body")
		return
	}
	defer r.Body.Close()

	if err := h.services.Users.SignUp(r.Context(), service.SignUpInput{
		LastName:  inp.LastName,
		FirstName: inp.FirstName,
		Email:     inp.Email,
		Password:  inp.Password,
	}); err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			response.NewResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		response.NewResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var inp service.SignInInput

	if err := json.NewDecoder(r.Body).Decode(&inp); err != nil {
		response.NewResponse(w, http.StatusBadRequest, "invalid input body")
		return
	}
	defer r.Body.Close()

	tokens, err := h.services.Users.SignIn(r.Context(), service.SignInInput{
		Email:    inp.Email,
		Password: inp.Password,
	})
	if err != nil {
		if err != nil {
			if errors.Is(err, domain.ErrUserNotFound) {
				response.NewResponse(w, http.StatusNotFound, err.Error())
				return
			}
			response.NewResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	res, err := json.Marshal(tokens)
	if err != nil {
		response.NewResponse(w, http.StatusInternalServerError, err.Error())
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
