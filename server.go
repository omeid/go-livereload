package livereload

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type nothing struct{}

type Server struct {
	Broadcast chan<- interface{}

	lock        sync.Mutex
	connections map[*websocket.Conn]nothing
}

func (s *Server) run(broadcast <-chan interface{}) {
	for m := range broadcast {
		s.lock.Lock()
		for conn, _ := range s.connections {
			if conn == nil {
				continue
			}
			err := conn.WriteJSON(m)
			if err != nil {
				Log.Println(err)
				delete(s.connections, conn)
			}
		}
		s.lock.Unlock()
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	for conn, _ := range s.connections {
		delete(s.connections, conn)
		conn.Close()
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		Log.Println(err)
		http.Error(w, "Can't upgrade.", 500)
		return
	}

	hello := helloMessage
	err = conn.ReadJSON(&hello)
	if err != nil {
		Log.Println(err)
		conn.Close()
		return
	}
	//TODO: Check protos compat.
	err = conn.WriteJSON(&serverHello)
	if err != nil {
		Log.Println(err)
		conn.Close()
		return
	}

	s.lock.Lock()
	s.connections[conn] = nothing{}
	s.lock.Unlock()
}

func (s *Server) Close() {
	close(s.Broadcast)
}

func (s *Server) Update(url string) {
	u := updateMessage
	u.Url = url
	s.Broadcast <- u
}

func (s *Server) Reload(path string) {
	r := reloadMessage
	r.Path = path
	s.Broadcast <- r
}

func (s *Server) Alert(alert string) {
	a := alertMessage
	a.Message = alert
	s.Broadcast <- a
}
