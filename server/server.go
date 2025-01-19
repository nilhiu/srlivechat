package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

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

func (s *Server) Run(ctx context.Context) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := s.sock.HandleRequest(w, r); err != nil {
			errMsg, _ := json.Marshal(
				newServerMessage(
					"ERROR",
					fmt.Sprintf("the broadcast server has encountered an error, %s", err.Error()),
				),
			)
			_ = s.sock.Broadcast(errMsg)
		}
	})

	s.sock.HandleMessage(func(session *melody.Session, msg []byte) {
		log.Info().Msgf("broadcasting message: %s", msg[:len(msg)-1])
		_ = s.sock.Broadcast(msg)
	})

	go func() {
		<-ctx.Done()
		shutDownMsg, _ := json.Marshal(
			newServerMessage("SHUTDOWN", "the server is shutting down..."),
		)

		_ = s.sock.Broadcast(shutDownMsg)
		s.sock.Close()

		log.Info().Msg("interrupt detected, ending server session...")
		time.Sleep(100 * time.Millisecond)
		os.Exit(0)
	}()

	return http.ListenAndServe(s.addr, nil)
}
