package workspace

import (
	"net/http"

	"github.com/Mikhail-tal63/viltrum_empier/middleware"
	"github.com/Mikhail-tal63/viltrum_empier/utils/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WorkspaceHandler struct {
	WorkspaceService *WorkspaceService
}

func NewWorkspaceHandler(WorkspaceService *WorkspaceService) *WorkspaceHandler {
	return &WorkspaceHandler{
		WorkspaceService: WorkspaceService,
	}

}

func (h *WorkspaceHandler) WorkspaceRouter(router *mux.Router) {
	router.HandleFunc("/workspaces", h.CreateWorkspace).Methods("POST")
	router.HandleFunc("/workspaces", h.ListWorkspaces).Methods("GET")
	router.HandleFunc("/workspaces/{workspace_id}/invite", h.InviteUser).Methods("POST")
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
	invitedUserID, err := primitive.ObjectIDFromHex(payload.InvitedUserID)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
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
	invite, err := h.WorkspaceService.InviteUser(
		r.Context(),
		userID,
		invitedUserID,
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
