package tui

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/fatih/color"
	"github.com/nilhiu/srlivechat/client"
	"github.com/nilhiu/srlivechat/server"
	"github.com/rs/zerolog/log"
)

var (
	colorUser       = color.New(color.Bold, color.FgCyan)
	colorConnect    = color.New(color.Bold, color.FgGreen)
	colorDisconnect = color.New(color.Bold, color.FgRed)
	colorServer     = color.New(color.Bold, color.FgMagenta)
)

type TUI interface {
	Run()
}

func New(ctx context.Context, c *client.Client) TUI {
	return &tui{
		ctx: ctx,
		c:   c,
	}
}

type tui struct {
	ctx context.Context
	c   *client.Client
}

func (t *tui) Run() {
	go t.shutdownHandler()
	go t.messageHandler()

	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			log.Error().Msgf("could not read written message, %s", err.Error())
		}

		if err := t.c.Write(scanner.Text()); err != nil {
			panic(err)
		}
	}
}

func (t *tui) shutdownHandler() {
	<-t.ctx.Done()
	log.Info().Msg("interrupt detected, ending client session...")
	t.c.Close()
	os.Exit(0)
}

func (t *tui) messageHandler() {
	for {
		msg, err := t.c.Read()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			log.Fatal().Msgf("failed to read from server, %s", err.Error())
		}

		switch msg.Type() {
		case server.UserMessage:
			t.userMessageHandler(msg)
		case server.ServerMessage:
			t.serverMessageHandler(msg)
		case server.ConnectMessage:
			t.connectMessageHandler(msg)
		case server.DisconnectMessage:
			t.disconnectMessageHandler(msg)
		}
	}
}

func (t *tui) userMessageHandler(msg server.Message) {
	userText := colorUser.Sprintf("[%s]:", msg.Sender())
	fmt.Printf("%s %s\n", userText, msg.Message())
}

func (t *tui) serverMessageHandler(msg server.Message) {
	svrText := colorServer.Sprint("<SERVER>:")
	fmt.Printf("%s %s\n", svrText, msg.Message())
	switch msg.Sender() {
	case "SHUTDOWN":
		os.Exit(0)
	case "CONFLICT":
		os.Exit(1)
	}
}

func (t *tui) connectMessageHandler(msg server.Message) {
	connText := colorConnect.Sprint("<CONNECTED>:")
	fmt.Printf("%s %s\n", connText, msg.Sender())
}

func (t *tui) disconnectMessageHandler(msg server.Message) {
	disconnText := colorDisconnect.Sprint("<DISCONNECTED>:")
	fmt.Printf("%s %s\n", disconnText, msg.Sender())
}
