package server

import "github.com/gorilla/websocket"

type MessageType uint

const (
	UserMessage = iota
	ServerMessage
	ConnectMessage
	DisconnectMessage
)

type Message interface {
	Type() MessageType
	Sender() string
	Message() string
}

type internalMessage struct {
	MsgType MessageType `json:"type"`
	User    string      `json:"name"`
	Msg     string      `json:"message"`
}

func (i *internalMessage) Type() MessageType {
	return i.MsgType
}

func (i *internalMessage) Sender() string {
	return i.User
}

func (i *internalMessage) Message() string {
	return i.Msg
}

func ReadMessage(conn *websocket.Conn) (Message, error) {
	var msg internalMessage
	if err := conn.ReadJSON(&msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

func NewUserMessage(user, message string) Message {
	return &internalMessage{
		MsgType: UserMessage,
		User:    user,
		Msg:     message,
	}
}

func NewConnectMessage(user string) Message {
	return &internalMessage{
		MsgType: ConnectMessage,
		User:    user,
		Msg:     "",
	}
}

func NewDisconnectMessage(user string) Message {
	return &internalMessage{
		MsgType: DisconnectMessage,
		User:    user,
		Msg:     "",
	}
}

func newServerMessage(cause, message string) Message {
	return &internalMessage{
		MsgType: ServerMessage,
		User:    cause,
		Msg:     message,
	}
}
