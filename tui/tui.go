package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nilhiu/srlivechat/client"
	"github.com/nilhiu/srlivechat/server"
	"github.com/rs/zerolog/log"
)

type TUI interface {
	Run()
}

func New(ctx context.Context, c *client.Client) TUI {
	return &tui{
		ctx: ctx,
		c:   c,
		p:   tea.NewProgram(initialModel(c), tea.WithAltScreen()),
	}
}

type tui struct {
	ctx context.Context
	c   *client.Client
	p   *tea.Program
}

func (t *tui) Run() {
	go func() {
		for {
			msg, err := t.c.Read()
			if err != nil {
				t.p.Send(err)
				return
			}

			t.p.Send(msg)
		}
	}()

	if _, err := t.p.Run(); err != nil {
		log.Fatal().Err(err)
	}
}

type model struct {
	c    *client.Client
	vp   viewport.Model
	msgs []string
	ti   textinput.Model

	sentStyle lipgloss.Style
	recvStyle lipgloss.Style
}

func initialModel(c *client.Client) model {
	tea.WindowSize()

	ti := textinput.New()
	ti.Placeholder = "Type a message..."
	ti.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Render("> ")
	ti.CharLimit = 64
	ti.Width = 30
	ti.Focus()

	vp := viewport.New(30, 10)
	vp.KeyMap = viewport.KeyMap{}

	return model{
		c:    c,
		vp:   vp,
		msgs: []string{},
		ti:   ti,

		sentStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("11")),
		recvStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("10")),
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.ti, tiCmd = m.ti.Update(msg)
	m.vp, vpCmd = m.vp.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.vp.Width = msg.Width
		m.ti.Width = msg.Width - 5
		m.vp.Height = msg.Height - 3

		if len(m.msgs) > 0 {
			m.vp.SetContent(strings.Join(m.msgs, "\n"))
		}
		m.vp.GotoBottom()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			m.c.Write(m.ti.Value())
			m.ti.Reset()
			m.vp.GotoBottom()
		}
	case server.Message:
		m.messageHandler(msg)
	case error:
		m.errorHandler(msg)
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	tiView := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("3")).
		Render(m.ti.View())

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.vp.View(),
		tiView,
	)
}

func (m *model) errorHandler(err error) {
	m.msgs = append(
		m.msgs,
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("0")).
			Background(lipgloss.Color("1")).
			Width(m.vp.Width).
			Render("<ERROR>: "+err.Error()),
	)

	m.ti.Reset()
	m.ti.Placeholder = "Not able to send messages, please restart."
	m.ti.Blur()

	m.vp.SetContent(strings.Join(m.msgs, "\n"))
	m.vp.GotoBottom()
}

func (m *model) messageHandler(msg server.Message) {
	switch msg.Type() {
	case server.ConnectMessage:
		m.msgs = append(
			m.msgs,
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("2")).
				Render("<CONNECTED>: "+msg.Sender()),
		)
	case server.DisconnectMessage:
		m.msgs = append(
			m.msgs,
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("1")).
				Render("<DISCONNECTED>: "+msg.Sender()),
		)
	case server.UserMessage:
		m.userMessageHandler(msg)
	case server.ServerMessage:
		m.msgs = append(
			m.msgs,
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("0")).
				Background(lipgloss.Color("5")).
				Width(m.vp.Width).
				Render("<SERVER>: "+msg.Message()),
		)
	}

	m.vp.SetContent(strings.Join(m.msgs, "\n"))
	m.vp.GotoBottom()
}

func (m *model) userMessageHandler(msg server.Message) {
	if msg.Sender() == m.c.Name() {
		m.msgs = append(
			m.msgs,
			fmt.Sprintf("%s: %s", m.sentStyle.Render(msg.Sender()), msg.Message()),
		)
	} else {
		m.msgs = append(
			m.msgs,
			fmt.Sprintf("%s: %s", m.recvStyle.Render(msg.Sender()), msg.Message()),
		)
	}
}
