package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/olahol/melody"
	"github.com/rs/zerolog/log"
)

type Server struct {
	sock  *melody.Melody
	addr  string
	users []string
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
		var m internalMessage
		if err := json.Unmarshal(msg, &m); err != nil {
			log.Error().Msgf("failed to read message: %s, (%s)", msg[:len(msg)-1], err.Error())
			return
		}

		switch m.MsgType {
		case ConnectMessage:
			if slices.Contains(s.users, m.User) {
				cMsg, _ := json.Marshal(
					newServerMessage("CONFLICT", "user with that username already in chat"),
				)
				session.Write(cMsg)
				time.Sleep(time.Second)
				session.Close()
				return
			}
			s.users = append(s.users, m.User)
		case DisconnectMessage:
			s.users = slices.DeleteFunc(s.users, func(u string) bool { return u == m.User })
		}

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

		for i := 3; i > 0; i-- {
			log.Info().Msgf("interrupt detected, ending server session in %d...", i)
			time.Sleep(time.Second)
		}
		os.Exit(0)
	}()

	return http.ListenAndServe(s.addr, nil)
}
