package controller

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"net/http"
	userRoleEnum "relay-backend/internal/enums"
	"relay-backend/internal/enums/sender"
	"relay-backend/internal/service"
	"relay-backend/internal/store"
	"relay-backend/internal/utils/exception"
)

type ChatController struct {
	chatService    *service.ChatService
	messageService *service.MessageService
	middleware     *AuthMiddleware
}

var (
	clientMap         = make(map[string]*websocket.Conn)
	operatorMap       = make(map[string]*websocket.Conn)
	operatorOnlineMap = make(map[int]*websocket.Conn)
)

type message struct {
	Sender string `json:"sender"`
	Text   string `json:"text"`
}

var senderMap = map[string]sender.Type{
	"system":   sender.System,
	"operator": sender.Operator,
	"client":   sender.Client,
}

func NewChatController(s *store.Store, authMiddleware *AuthMiddleware) func(r chi.Router) {
	cc := &ChatController{
		chatService:    service.NewChatService(s),
		messageService: service.NewMessageService(s),
		middleware:     authMiddleware,
	}

	return func(r chi.Router) {
		r.Get("/ws", cc.wsHandler)
		r.Post("/{widgetUuid}", cc.createChat)
		r.With(authMiddleware.Auth(userRoleEnum.Access.Operator)).Get("/{orgId}/list", cc.wsOperatorChatListHandler)
	}
}

func (cc *ChatController) createChat(w http.ResponseWriter, r *http.Request) {
	widgetUuid := chi.URLParam(r, "widgetUuid")
	chat, err := cc.chatService.CreateNewChat(widgetUuid)
	if err != nil {
		HTTPError(w, r, err.(exception.Exception))
		return
	}

	for _, conn := range operatorOnlineMap {
		conn.WriteMessage(websocket.TextMessage, []byte(chat.Uuid))
	}

	Respond(w, r, http.StatusOK, &chat)
}

func (cc *ChatController) wsHandler(w http.ResponseWriter, r *http.Request) {
	chatId := r.URL.Query().Get("chatId")
	sndr := r.URL.Query().Get("sender")

	from := senderMap[sndr]

	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		HTTPError(w, r, exception.NewDetailsException(http.StatusInternalServerError, exception.Enum.InternalServerError, map[string]interface{}{"error": err}))
		return
	}

	defer conn.Close()

	//TODO: before ws connection widget must send request to create chat, which receives cookie with chat.id
	//TODO: get chat by cookie
	//TODO: if there is no cookie then close ws connection

	//TODO: consider about endpoint with auth for operators and THEN add them to map
	//consider (need to AUTHORIZE operator) TODO: authorize client (by widget uuid and cookie) or operator (by user.id and organization.id) and send all previous messages
	//TODO: need to check chat existing
	if from == sender.Client {
		if _, ok := clientMap[chatId]; !ok {
			clientMap[chatId] = conn
		}
	} else if from == sender.Operator {
		if _, ok := operatorMap[chatId]; !ok {
			operatorMap[chatId] = conn
		}
	}

	for {
		var msg message

		//TODO: if there is no receiver then do nothing (now is err and ws disconnecting)
		//TODO: check if operator online, if not then response about waiting

		messageType, m, err := conn.ReadMessage()
		if err != nil {
			HTTPError(w, r, exception.NewDetailsException(http.StatusInternalServerError, exception.Enum.InternalServerError, map[string]interface{}{"error": err}))
			return
		}

		if err := json.Unmarshal(m, &msg); err != nil {
			HTTPError(w, r, exception.NewDetailsException(http.StatusInternalServerError, exception.Enum.InternalServerError, map[string]interface{}{"error": err}))
			return
		}

		msgSender := senderMap[msg.Sender]
		if msgSender == sender.Operator {
			if client, ok := clientMap[chatId]; ok {
				err := client.WriteMessage(messageType, []byte(msg.Text))
				if err != nil {
					return
				}
			}
		} else if msgSender == sender.Client {
			if operator, ok := operatorMap[chatId]; ok {
				err := operator.WriteMessage(messageType, []byte(msg.Text))
				if err != nil {
					return
				}
			}
		}

		err = cc.messageService.SaveMessage(msgSender, msg.Text, 0, chatId)
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		}
	}
}

func (cc *ChatController) wsOperatorChatListHandler(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(CtxKeyUser).(int)

	upgrader := websocket.Upgrader{}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		HTTPError(w, r, exception.NewDetailsException(http.StatusInternalServerError, exception.Enum.InternalServerError, map[string]interface{}{"err": err}))
		return
	}

	defer func() {
		conn.Close()
		delete(operatorOnlineMap, id)
	}()

	operatorOnlineMap[id] = conn

	//TODO: notification for others operators that chat is already serving
	for {
		messageType, m, err := conn.ReadMessage()
		if err != nil {
			HTTPError(w, r, exception.NewDetailsException(http.StatusInternalServerError, exception.Enum.InternalServerError, map[string]interface{}{"error": err}))
			return
		}

		fmt.Println(m)
		conn.WriteMessage(messageType, m)
	}
}
