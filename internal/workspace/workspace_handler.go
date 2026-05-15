package workspace

import (
	"net/http"

	"github.com/Mikhail-tal63/viltrum_empier/middleware"
	"github.com/Mikhail-tal63/viltrum_empier/utils/json"
	"github.com/gorilla/mux"
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
