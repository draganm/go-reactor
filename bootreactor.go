package reactor

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
)

// UserEvent is an event triggered by the client. Such as click, or key events.
type UserEvent struct {
	ElementID   string                 `json:"id,omitempty"`
	Type        string                 `json:"type,omitempty"`
	Value       string                 `json:"value,omitempty"`
	Data        string                 `json:"data,omitempty"`
	ExtraValues map[string]interface{} `json:"xv,omitempty"`
}

// ClientConnectionHandler interacts with a newly connected user.
type ClientConnectionHandler func(chan *DisplayUpdate, chan *UserEvent, *http.Request, string) http.Header

// newReactorHandler creates a new HTTP Handler for user WebsSocket connections.
func newReactorHandler(listener ClientConnectionHandler) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		displayChan := make(chan *DisplayUpdate)
		eventChan := make(chan *UserEvent)

		connectionID := uuid.NewV4().String()

		header := listener(displayChan, eventChan, r, connectionID)

		conn, err := upgrader.Upgrade(w, r, header)
		if err != nil {
			panic(err)
		}

		go func() {
			currentState := &DisplayUpdate{}
			for displayUpdate := range displayChan {
				if !displayUpdate.DeepEqual(currentState) {
					err := conn.WriteJSON(displayUpdate)
					if err != nil {
						break
					}
					currentState = displayUpdate
				}
			}
		}()

		defer func() {
			close(eventChan)
		}()

		for {
			evt := &UserEvent{}
			err := conn.ReadJSON(evt)
			if err != nil {
				// fmt.Println("should close! ", err)
				return
			}

			eventChan <- evt

		}

	}

}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}
