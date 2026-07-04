package messages

import (
	"net/http"

	"github.com/Mikhail-tal63/viltrum_empier/middleware"
	"github.com/Mikhail-tal63/viltrum_empier/utils/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageHanler struct{
	service *MessageService
}

func NewMessageHandler(service *MessageService)*MessageHanler{
	return &MessageHanler{
		service: service,
	}
}

func (h *MessageHanler)MessageRouter(router *mux.Router){
router.HandleFunc("/group-chat/{groupChatId}/messages",h.ListMessagesInGroupChat).Methods("GET")
router.HandleFunc("/group-chat/{groupChatId}/message",h.CreateMessage).Methods("POST")
}

func (h *MessageHanler) ListMessagesInGroupChat(w http.ResponseWriter,r *http.Request){
	
	gID,err := primitive.ObjectIDFromHex(mux.Vars(r)["groupChatId"])
	if err != nil {
		json.WriteError(w,http.StatusInternalServerError,err)
		return
	}
	messages, err := h.service.ListMessagesInGroupChat(r.Context(),gID)
if err != nil {
		json.WriteError(w,http.StatusInternalServerError,err)
		return
}
if err := json.WriteJSON(w,http.StatusOK,messages)
err != nil {
		json.WriteError(w,http.StatusInternalServerError,err)
		return
}
}


func (h *MessageHanler) CreateMessage(w http.ResponseWriter,r *http.Request){
	var payload CreateMessagePayload
	if err := json.ParseJSON(r,&payload)
	err != nil {
		json.WriteError(w,http.StatusInternalServerError,err)
		return
	}
	userid ,err := middleware.GetUserID(r.Context())
	if err != nil {
		json.WriteError(w,http.StatusInternalServerError,err)
		return
	}
	gID,err := primitive.ObjectIDFromHex(mux.Vars(r)["groupChatId"])
	if err != nil {
		json.WriteError(w,http.StatusInternalServerError,err)
		return
	}
	msg ,err := h.service.CreateMessage(r.Context(),&payload,userid,gID)
	if err != nil {
		json.WriteError(w,http.StatusInternalServerError,err)
		return
	}
	if err := json.WriteJSON(w,http.StatusOK,msg)
	err != nil {
	json.WriteError(w,http.StatusInternalServerError,err)
		return
	}
}
