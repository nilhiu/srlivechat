package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/olahol/melody"
	"github.com/rs/zerolog/log"
)

type Server struct {
	sock *melody.Melody
	addr string
}

func New(addr string) *Server {
	return &Server{
		sock: melody.New(),
		addr: addr,
	}
}

func (s *Server) Run() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := s.sock.HandleRequest(w, r); err != nil {
			errBroadcast, _ := json.Marshal(Message{
				User: "SERVER",
				Message: fmt.Sprintf(
					"the broadcast server has encountered an error: %s",
					err.Error(),
				),
			})
			_ = s.sock.Broadcast(errBroadcast)
		}
	})

	s.sock.HandleConnect(func(session *melody.Session) {
		log.Info().Msg("user connected")
	})

	s.sock.HandleMessage(func(session *melody.Session, msg []byte) {
		log.Info().Msgf("broadcasting message: %s", msg[:len(msg)-1])
		_ = s.sock.Broadcast(msg)
	})

	return http.ListenAndServe(s.addr, nil)
}
