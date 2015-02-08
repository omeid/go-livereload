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

	host   = "devel"
	protos = []string{
		"http://livereload.com/protocols/official-7",
		"http://livereload.com/protocols/official-8",
		"http://livereload.com/protocols/official-9",
		"http://livereload.com/protocols/2.x-origin-version-negotiation",
		"http://livereload.com/protocols/2.x-remote-control",
	}
)

func New() *Server {
	broadcast := make(chan interface{}, 50)
	s := &Server{
		Broadcast:   broadcast,
		connections: make(map[*websocket.Conn]nothing),
	}
	go s.run(broadcast)
	return s
}
