package task

import (
	"net/http"

	"github.com/Mikhail-tal63/viltrum_empier/middleware"
	"github.com/Mikhail-tal63/viltrum_empier/utils/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskHandler struct {
	service *TaskService
}

func NewTaskHandler(service *TaskService) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}
func (h *TaskHandler) TaskRouter(router *mux.Router) {
	router.HandleFunc("/columns/{column_id}/tasks", h.CreateTask).Methods("POST")
	router.HandleFunc("/columns/{column_id}/tasks", h.ListTasks).Methods("GET")
}
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var payload CreateTaskPayload

	if err := json.ParseJSON(r, &payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	columnID, err := primitive.ObjectIDFromHex(mux.Vars(r)["column_id"])
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
	}
	task, err := h.service.CreateTask(r.Context(), &payload, userID, columnID)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := json.WriteJSON(w, http.StatusCreated, task); err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	columnID, err := primitive.ObjectIDFromHex(mux.Vars(r)["column_id"])
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}

	tasks, err := h.service.ListTasks(r.Context(), columnID)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := json.WriteJSON(w, http.StatusOK, tasks); err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

}
