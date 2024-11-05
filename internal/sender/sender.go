package sender

import (
	"context"
	"errors"

	"notify/internal/parser"
)

type Sender interface {
	Send(ctx context.Context, content string, summary string, extra map[string]any) error
}

type Manager struct {
	senders map[parser.Platform]Sender
}

func NewManager() *Manager {
	return &Manager{
		senders: make(map[parser.Platform]Sender),
	}
}

func (m *Manager) Register(platform parser.Platform, sender Sender) {
	m.senders[platform] = sender
}

func (m *Manager) Send(ctx context.Context, msg *parser.Message) error {
	sender, ok := m.senders[msg.Platform]
	if !ok {
		return errors.New("unsupported platform")
	}

	return sender.Send(ctx, msg.Content, msg.Summary, msg.Extra)
}
