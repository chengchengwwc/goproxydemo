package server

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type webSocketHandler struct {
	upgrader       *websocket.Upgrader
	binder         *binder
	calcUserIDFunc func(token string) (userID string, ok bool)
}

type RegisterMessage struct {
	Token string
	Event string
}

func (wb *webSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wsConn, err := wb.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer wsConn.Close()
}
