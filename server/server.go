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
	ctx   context.Context
	sock  *melody.Melody
	addr  string
	users []string
}

func New(ctx context.Context, addr string) *Server {
	return &Server{
		ctx:  ctx,
		sock: melody.New(),
		addr: addr,
	}
}

func (s *Server) Run() error {
	http.HandleFunc("/", s.WebSocketHandler)
	s.sock.HandleMessage(s.MessageHandler)
	go s.ShutdownHandler()

	return http.ListenAndServe(s.addr, nil)
}

func (s *Server) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	if err := s.sock.HandleRequest(w, r); err != nil {
		errMsg, _ := json.Marshal(
			newServerMessage(
				"ERROR",
				fmt.Sprintf("the broadcast server has encountered an error, %s", err.Error()),
			),
		)
		_ = s.sock.Broadcast(errMsg)
	}
}

func (s *Server) MessageHandler(session *melody.Session, msg []byte) {
	var m internalMessage
	if err := json.Unmarshal(msg, &m); err != nil {
		log.Error().Msgf("failed to read message: %s, (%s)", msg[:len(msg)-1], err.Error())
		return
	}

	switch m.MsgType {
	case ConnectMessage:
		s.HandleConnection(session, m.User)
	case DisconnectMessage:
		s.HandleDisconnection(session, m.User)
	case UserMessage:
		s.sock.Broadcast(msg)
	}

	log.Info().Msgf("broadcasting message: %s", msg[:len(msg)-1])
}

func (s *Server) ShutdownHandler() {
	<-s.ctx.Done()
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
}

func (s *Server) HandleConnection(session *melody.Session, user string) {
	if slices.Contains(s.users, user) {
		cMsg, _ := json.Marshal(
			newServerMessage("CONFLICT", "user with that username already in chat"),
		)
		session.Write(cMsg)
		time.Sleep(time.Second)
		session.Close()
		return
	}

	s.users = append(s.users, user)
}

func (s *Server) HandleDisconnection(session *melody.Session, user string) {
	s.users = slices.DeleteFunc(s.users, func(u string) bool {
		return u == user
	})
}
