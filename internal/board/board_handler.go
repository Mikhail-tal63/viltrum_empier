package board

import (
	"net/http"

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
	router.HandleFunc("/boards", h.CreateBoard).Methods("POST")
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
	board, err := h.service.CreateBoard(r.Context(), &payload, workspaceID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.WriteJSON(w, http.StatusCreated, board); err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
	}

}
