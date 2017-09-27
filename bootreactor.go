package reactor

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// UserEvent is an event triggered by the client. Such as click, or key events.
type UserEvent struct {
	ElementID   string                 `json:"id,omitempty"`
	Type        string                 `json:"type,omitempty"`
	Value       string                 `json:"value,omitempty"`
	Data        string                 `json:"data,omitempty"`
	ExtraValues map[string]interface{} `json:"xv,omitempty"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}
