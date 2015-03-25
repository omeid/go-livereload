package livereload

import (
	"log"
	"os"

	"github.com/gorilla/websocket"
)

type Logger interface {
	Println(...interface{})
}

var (
	Log Logger = log.New(os.Stderr, "", log.LstdFlags)

	protos = []string{
		"http://livereload.com/protocols/official-7",
	}
)

func New(name string) *Server {
	broadcast := make(chan interface{}, 50)
	s := &Server{
	  name: name,
		broadcast:   broadcast,
		connections: make(map[*websocket.Conn]nothing),
	}
	go s.run(broadcast)
	return s
}
