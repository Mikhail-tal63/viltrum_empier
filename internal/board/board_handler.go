package board

import (
	"net/http"

	"github.com/Mikhail-tal63/viltrum_empier/middleware"
	"github.com/Mikhail-tal63/viltrum_empier/utils/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BoardHandler struct {
	service *BoardService
}

func NewBoardHandler(service *BoardService) *BoardHandler {
	return &BoardHandler{
		service: service,
	}
}

func (h *BoardHandler) BoardRouter(router *mux.Router) {
	router.HandleFunc("/workspaces/{workspace_id}/boards", h.CreateBoard).Methods("POST")
	router.HandleFunc("/workspaces/{workspace_id}/boards", h.GetWorkspaceBoards).Methods("GET")
	router.HandleFunc("/boards/{board_id}/colmuns", h.CreateColmun).Methods("POST")
	router.HandleFunc("/boards/{board_id}/columns", h.GetBoarderColmuns).Methods("GET")

}

func (h *BoardHandler) CreateBoard(w http.ResponseWriter, r *http.Request) {
	var payload CreateBoardPayload

	if err := json.ParseJSON(r, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	workspaceID, err := primitive.ObjectIDFromHex(mux.Vars(r)["workspace_id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}
	userid, err := middleware.GetUserID(r.Context())
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	board, err := h.service.CreateBoard(r.Context(), &payload, userid, workspaceID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.WriteJSON(w, http.StatusCreated, board); err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)

	}

}

func (h *BoardHandler) GetWorkspaceBoards(w http.ResponseWriter, r *http.Request) {
	workspaceid, err := primitive.ObjectIDFromHex(mux.Vars(r)["workspace_id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	boards, err := h.service.GetWorkspaceBoards(r.Context(), workspaceid)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if err := json.WriteJSON(w, http.StatusOK, boards); err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
	}

}

func (h *BoardHandler) DeleteBoard(w http.ResponseWriter, r *http.Request) {
	boardID, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := h.service.DeleteBoard(r.Context(), userID, boardID); err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := json.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "board deleted",
	}); err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *BoardHandler) CreateColmun(w http.ResponseWriter, r *http.Request) {
	var payload ColumnPayload

	if err := json.ParseJSON(r, &payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	boardID, err := primitive.ObjectIDFromHex(mux.Vars(r)["board_id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userID, err := middleware.GetUserID(r.Context())
	colmun, err := h.service.CreateColmun(r.Context(), boardID, userID, &payload)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if err := json.WriteJSON(w, http.StatusCreated, colmun); err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
	}
}

func (h *BoardHandler) GetBoarderColmuns(w http.ResponseWriter, r *http.Request) {
	boardID, err := primitive.ObjectIDFromHex(mux.Vars(r)["board_id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	colmuns, err := h.service.GetBoarderColumns(r.Context(), boardID)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if err := json.WriteJSON(w, http.StatusOK, colmuns); err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
	}

}

func (h *BoardHandler) DeleteColumn(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	userid, err := middleware.GetUserID(r.Context())
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if err := h.service.DeleteColumn(r.Context(), id, userid); err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}
