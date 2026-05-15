package auth

import (
	"net/http"

	"github.com/Mikhail-tal63/viltrum_empier/utils/json"
	"github.com/gorilla/mux"
)

type AuthHandler struct {
	service *AuthService
}

func NewAuthHandler(service *AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) AuthRouter(router *mux.Router) {
	router.HandleFunc("/register", h.CreateUser).Methods("POST")
	router.HandleFunc("/login", h.Login).Methods("POST")
}

func (h *AuthHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var payload CreateUserPayload
	if err := json.ParseJSON(r, &payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}
	crerateduser, err := h.service.CreateUser(&payload)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := json.WriteJSON(w, http.StatusCreated, crerateduser); err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return

	}

}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var payload LoginPayload
	if err := json.ParseJSON(r, &payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}
	user, err := h.service.Login(payload.Email, payload.Password)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if err := json.WriteJSON(w, http.StatusOK, user); err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return

	}
}
