package ws

type hub struct {
	// broadcast     chan []byte
	// room          chan roomInfo
	// broadcastNoMe chan broadcastNoMeInfo
	clients    map[*Socket]bool
	register   chan *Socket
	unregister chan *Socket
}

func newHub() *hub {
	return &hub{
		register:   make(chan *Socket),
		unregister: make(chan *Socket),
		clients:    make(map[*Socket]bool),
	}
}

func (h *hub) Run() {
	for {
		select {
		case socket := <-h.register:
			h.clients[socket] = true
		case socket := <-h.unregister:
			delete(h.clients, socket)
			// close(client.send)
			// case message := <-h.broadcast:
			// 	for id, client := range h.clients {
			// 		select {
			// 		case client.send <- message:
			// 		default:
			// 			delete(h.clients, id)
			// 			close(client.send)
			// 		}
			// 	}
			// case info := <-h.room:
			// 	for id, client := range h.clients {
			// 		for _, room := range info.rooms {
			// 			_, ok := client.rooms[room]
			// 			// log.Println(client.id, client.rooms)
			// 			if ok {
			// 				select {
			// 				case client.send <- info.message:
			// 				default:
			// 					delete(h.clients, id)
			// 					close(client.send)
			// 				}
			// 			}
			// 		}

			// 	}
			// case info := <-h.broadcastNoMe:
			// 	for id, client := range h.clients {
			// 		if info.clientID != client.id {
			// 			select {
			// 			case client.send <- info.message:
			// 			default:
			// 				delete(h.clients, id)
			// 				close(client.send)
			// 			}
			// 		}
			// 	}
		}
	}
}

// type Server struct {
// 	// hub        *hub
// 	// onCallback map[string]func(msg json.RawMessage)
// }

// func newServer(hub *hub) *Server {
// 	return &Server{
// 		// hub:        hub,
// 		onCallback: make(map[string]func(msg json.RawMessage)),
// 	}
// }

// type roomInfo struct {
// 	rooms   []string
// 	message []byte
// }

// type broadcastNoMeInfo struct {
// 	clientID string
// 	message  []byte
// }

// type Message struct {
// 	Message string `json:"message"`
// }

// func T() {
// 	New(func(socket *Socket) {
// 		socket.On("message", func(msg []byte) {
// 			var m Message
// 			if err := json.Unmarshal(msg, &m); err != nil {
// 				log.Println(err)
// 			}
// 			log.Println(m)
// 		})

// 	})
// }
