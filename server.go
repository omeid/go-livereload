package livereload

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type nothing struct{}

type Server struct {
	broadcast chan<- interface{}

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
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var sHello = serverHello{message{"hello"}, protos, host}

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

	chello := helloMessage{}
	err = conn.ReadJSON(&chello)
	if err != nil {
		Log.Println(err)
		conn.Close()
		return
	}
	//TODO: Check protos compat.
	err = conn.WriteJSON(&sHello)
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
	close(s.broadcast)
}

func (s *Server) Update(url string) {
	s.broadcast <- updateMessage{
		message{"update"},
		url,
	}
}

func (s *Server) Reload(path string, cssLivereload bool) {
	s.broadcast <- reloadMessage{
		message{"reload"},
		path,
		cssLivereload,	
	}
}

func (s *Server) Alert(alert string) {
	s.broadcast <- alertMessage{
		message{"alert"},
		alert,
	}
}
