package messages

import (
	"net/http"

	"github.com/Mikhail-tal63/viltrum_empier/middleware"
	"github.com/Mikhail-tal63/viltrum_empier/utils/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageHanler struct {
	service *MessageService
}

func NewMessageHandler(service *MessageService) *MessageHanler {
	return &MessageHanler{
		service: service,
	}
}

func (h *MessageHanler) MessageRouter(router *mux.Router) {
	router.HandleFunc("/group-chat/{groupChatId}/messages", h.ListMessagesInGroupChat).Methods("GET")
	router.HandleFunc("/group-chat/{groupChatId}/message", h.CreateMessage).Methods("POST")
    router.HandleFunc("/message/{id}",h.DeleteMessage).Methods("DELETE")
	router.HandleFunc("/message/{id}",h.EditMessage).Methods("PATCH")
}

func (h *MessageHanler) ListMessagesInGroupChat(w http.ResponseWriter, r *http.Request) {

	gID, err := primitive.ObjectIDFromHex(mux.Vars(r)["groupChatId"])
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}
	messages, err := h.service.ListMessagesInGroupChat(r.Context(), gID)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if err := json.WriteJSON(w, http.StatusOK, messages); err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *MessageHanler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var payload CreateMessagePayload
	if err := json.ParseJSON(r, &payload); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}
	userid, err := middleware.GetUserID(r.Context())
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	gID, err := primitive.ObjectIDFromHex(mux.Vars(r)["groupChatId"])
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}
	msg, err := h.service.CreateMessage(r.Context(), &payload, userid, gID)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if err := json.WriteJSON(w, http.StatusCreated, msg); err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *MessageHanler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	userid, err := middleware.GetUserID(r.Context())
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	message, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := h.service.DeleteMessage(r.Context(), message, userid); err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if err := json.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "message deleted seccsessfuly",
	},
	); err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *MessageHanler) EditMessage(w http.ResponseWriter,r *http.Request){

	var payload EditMessagePayload
 
	if err := json.ParseJSON(r,&payload);err != nil {
	
			json.WriteError(w, http.StatusBadRequest, err)
		return
	}
	userid, err := middleware.GetUserID(r.Context())
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	msgid,err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		json.WriteError(w,http.StatusBadRequest,err)
		return
	}
	if err := h.service.EditMessage(r.Context(),msgid,userid,&payload);err != nil {
	
		json.WriteError(w,http.StatusInternalServerError,err)
		return
	}
	if err := json.WriteJSON(w, http.StatusOK, map[string]string{
    "message": "message edited successfully",
}); err != nil {
    json.WriteError(w, http.StatusInternalServerError, err)
    return
}
}