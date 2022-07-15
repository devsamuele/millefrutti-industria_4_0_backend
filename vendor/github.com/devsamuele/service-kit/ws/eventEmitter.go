package ws

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/devsamuele/service-kit/web"
	"nhooyr.io/websocket"
)

type EventEmitter struct {
	hub    *hub
	config *Config
}

type Config websocket.AcceptOptions

func (e *EventEmitter) EmitTo(rooms []string, event string, msg json.RawMessage) error {
	for _, room := range rooms {
		for socket := range e.hub.clients {
			if _, ok := socket.rooms[room]; ok {
				im := internalMessage{
					Event: event,
					Data:  msg,
				}
				b, err := json.Marshal(&im)
				if err != nil {
					return err
				}
				socket.send <- b
			}
		}
	}
	return nil
}

func (e *EventEmitter) Broadcast(event string, msg json.RawMessage) error {
	for socket := range e.hub.clients {
		im := internalMessage{
			Event: event,
			Data:  msg,
		}
		b, err := json.Marshal(&im)
		if err != nil {
			return err
		}
		socket.send <- b
	}
	return nil
}

func New(cfg *Config) EventEmitter {
	ee := EventEmitter{
		hub:    newHub(),
		config: cfg,
	}
	go ee.hub.Run()
	return ee
}

func (ee *EventEmitter) OnConnection(callback func(r *http.Request, socket *Socket)) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
		if err != nil {
			return err
		}

		// id, err := uuid.NewRandom()
		// if err != nil {
		// 	return err
		// }

		s := Socket{
			conn:       conn,
			Meta:       make(map[string]string),
			send:       make(chan []byte),
			hub:        ee.hub,
			rooms:      make(map[string]bool),
			onCallback: make(map[string]func(msg json.RawMessage)),
		}
		ee.hub.register <- &s

		// log.Println("new client id", id)

		_ctx, cancel := context.WithCancel(context.Background())

		callback(r, &s)
		go s.read(_ctx, cancel, s.onCallback, s.onDisconnectCallback)
		go s.write(_ctx, cancel)

		return nil
	}
}
