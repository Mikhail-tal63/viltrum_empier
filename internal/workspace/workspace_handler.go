package workspace

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Mikhail-tal63/viltrum_empier/internal/auth"
	"github.com/Mikhail-tal63/viltrum_empier/middleware"
	"github.com/Mikhail-tal63/viltrum_empier/utils/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WorkspaceHandler struct {
	WorkspaceService *WorkspaceService
	AuthService      *auth.AuthService
}

func NewWorkspaceHandler(WorkspaceService *WorkspaceService, AuthService *auth.AuthService) *WorkspaceHandler {
	return &WorkspaceHandler{
		WorkspaceService: WorkspaceService,
		AuthService:      AuthService,
	}

}

func (h *WorkspaceHandler) WorkspaceRouter(router *mux.Router) {
	router.HandleFunc("/workspaces", h.CreateWorkspace).Methods("POST")
	router.HandleFunc("/workspaces", h.ListWorkspaces).Methods("GET")
	router.HandleFunc("/workspaces/invites", h.GetInviteByUserID).Methods("GET")
	router.HandleFunc("/workspaces/{workspace_id}/invite", h.InviteUser).Methods("POST")
	router.HandleFunc("/workspaces/{workspace_id}/invite/{invite_id}/accept", h.AcceptInvite).Methods("POST")
}

func (h *WorkspaceHandler) CreateWorkspace(w http.ResponseWriter, r *http.Request) {
	var payload WorkspacePayload
	if err := json.ParseJSON(r, &payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	workspace, err := h.WorkspaceService.CreateWorkspace(r.Context(), &payload, userID)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if err := json.WriteJSON(w, http.StatusCreated, workspace); err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *WorkspaceHandler) ListWorkspaces(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	Workspace, err := h.WorkspaceService.ListWorkspaces(r.Context(), userID)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if err := json.WriteJSON(w, http.StatusOK, Workspace); err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *WorkspaceHandler) InviteUser(w http.ResponseWriter, r *http.Request) {
	var payload InviteUserPayload

	if err := json.ParseJSON(r, &payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}
	username := strings.ToLower(strings.TrimSpace(payload.Username))
	if username == "" {
		json.WriteError(w, http.StatusBadRequest, fmt.Errorf("username is required"))
		return
	}
	invitedUser, err := h.AuthService.GetUserByUsername(r.Context(), username)
	if err != nil {
		json.WriteError(w, http.StatusNotFound, err)
		return
	}
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	workspaceID, err := primitive.ObjectIDFromHex(mux.Vars(r)["workspace_id"])
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if invitedUser.ID == userID {
		json.WriteError(w, http.StatusBadRequest, fmt.Errorf("you can't invite yourself"))
		return
	}
	invite, err := h.WorkspaceService.InviteUser(
		r.Context(),
		userID,
		invitedUser.ID,
		workspaceID,
	)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if err := json.WriteJSON(w, http.StatusCreated, invite); err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
	}
}

func (h *WorkspaceHandler) GetInviteByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	invites, err := h.WorkspaceService.GetInviteByUserID(r.Context(), userID)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := json.WriteJSON(w, http.StatusOK, invites); err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

}

func (h *WorkspaceHandler) AcceptInvite(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	inviteID, err := primitive.ObjectIDFromHex(mux.Vars(r)["invite_id"])
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}
	Workspace, err := primitive.ObjectIDFromHex(mux.Vars(r)["workspace_id"])
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}
	invite, err := h.WorkspaceService.AcceptInvite(r.Context(), Workspace, userID, inviteID)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if err := json.WriteJSON(w, http.StatusOK, invite); err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}
