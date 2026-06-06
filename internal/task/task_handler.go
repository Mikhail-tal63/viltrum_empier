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
	router.HandleFunc("/tasks/{_id}/dragdrop", h.DragDropTaskToColumnHandelr).Methods("POST")
	router.HandleFunc("/tasks/{_id}", h.EditTask).Methods("PUT")
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
		return
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

func (h *TaskHandler) DragDropTaskToColumnHandelr(w http.ResponseWriter, r *http.Request) {
	var payload DragPayload

	err := json.ParseJSON(r, &payload)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	taskID, err := primitive.ObjectIDFromHex(mux.Vars(r)["_id"])
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	newcolumnID, err := primitive.ObjectIDFromHex(payload.NewColumnID)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.service.DragDropTaskToColumn(r.Context(), payload.Position, taskID, userID, newcolumnID)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = json.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "task moved successfully",
	})
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

}

func (h *TaskHandler) EditTask(w http.ResponseWriter, r *http.Request) {
	var payload EditTaskPayload

	if err := json.ParseJSON(r, &payload); err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	taslID, err := primitive.ObjectIDFromHex(mux.Vars(r)["_id"])
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if err := h.service.EditTaskDetails(r.Context(), &payload, taslID, userID); err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := json.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "task edited successfully",
	}); err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
	}
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {

	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if err := h.service.DeleteTask(r.Context(), id, userID); err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := json.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "task deleted seccessfuly",
	}); err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}
