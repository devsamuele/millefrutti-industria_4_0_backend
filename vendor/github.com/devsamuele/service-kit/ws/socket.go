package ws

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"nhooyr.io/websocket"
)

// type client struct {
// 	conn  *websocket.Conn
// 	send  chan []byte
// 	rooms map[string]bool
// 	hub   *hub
// }

type Socket struct {
	Meta                 map[string]string
	conn                 *websocket.Conn
	send                 chan []byte
	rooms                map[string]bool
	hub                  *hub
	onCallback           map[string]func(msg json.RawMessage)
	onDisconnectCallback func()
}

// func (s Socket) GetID() string {
// 	return s.id
// }

func (s *Socket) On(event string, callback func(msg json.RawMessage)) {
	s.onCallback[event] = callback
	// getCallback
}

func (s *Socket) OnDisconnect(callback func()) {
	s.onDisconnectCallback = callback
	// getCallback
}

func (s *Socket) Join(rooms []string) {
	// log.Println("joined id:", s.client.id)
	for _, room := range rooms {
		s.rooms[room] = true
	}
}

func (s *Socket) Emit(event string, msg json.RawMessage) error {
	// s.hub.clients[s.id].send <- msg
	im := internalMessage{
		// SocketID: s.id,
		Event: event,
		Data:  msg,
	}
	b, err := json.Marshal(&im)
	if err != nil {
		return err
	}
	s.send <- b
	return nil
}

func (s *Socket) EmitTo(rooms []string, event string, msg json.RawMessage) error {
	for _, room := range rooms {
		for socket := range s.hub.clients {
			if s != socket {
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
	}
	return nil
}

func (s *Socket) Broadcast(event string, msg json.RawMessage) error {
	for socket := range s.hub.clients {
		if s != socket {
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
	return nil
}

type internalMessage struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}

func (s *Socket) read(ctx context.Context, cancel context.CancelFunc, eventFnMap map[string]func(msg json.RawMessage), onDisconnectCallback func()) {

	defer func() {
		s.hub.unregister <- s
		s.conn.Close(websocket.StatusGoingAway, "bye bye")
		if onDisconnectCallback != nil {
			onDisconnectCallback()
		}
		cancel()
	}()

	for {
		_, data, err := s.conn.Read(ctx)
		if err != nil {
			return
		}

		var im internalMessage
		if err := json.Unmarshal(data, &im); err != nil {
			log.Println("internal message error:", err)
			return
		}

		// log.Printf("\t EVENT: %v - ID: %v", im.Event, c.id)
		if f, ok := eventFnMap[im.Event]; ok {
			f(im.Data)
		}

	}
}

func (s *Socket) write(ctx context.Context, cancel context.CancelFunc) {
	defer func() {
		// c.hub.unregister <- c
		// log.Println("user disconnect")
		s.conn.Close(websocket.StatusGoingAway, "bye bye")
		cancel()
	}()

	ticker := time.NewTicker(time.Second * 60)

	for {
		select {
		case msg, ok := <-s.send:
			if !ok {
				// log.Println("close channel")
				return
			}

			if err := s.conn.Write(context.Background(), websocket.MessageText, msg); err != nil {
				log.Println(err)
				return
			}

		case <-ticker.C:
			if err := s.conn.Ping(ctx); err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
