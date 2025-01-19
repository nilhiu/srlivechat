package client

import (
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/nilhiu/srlivechat/server"
)

type Client struct {
	conn *websocket.Conn
	name string
}

func New(addr, name string) (*Client, error) {
	url := url.URL{
		Scheme: "ws",
		Host:   addr,
	}
	conn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		return nil, err
	}

	if err := conn.WriteJSON(server.NewConnectMessage(name)); err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
		name: name,
	}, nil
}

func (c *Client) Close() {
	_ = c.conn.WriteJSON(server.NewDisconnectMessage(c.name))
	c.conn.Close()
}

func (c *Client) Read() (server.Message, error) {
	msg, err := server.ReadMessage(c.conn)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (c *Client) Write(msg string) error {
	return c.conn.WriteJSON(server.NewUserMessage(c.name, msg))
}
