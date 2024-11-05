package dispatcher

import (
	"context"
	"sync"

	"notify/internal/parser"
	"notify/internal/sender"
	"notify/pkg/logger"

	"go.uber.org/zap"
)

type Dispatcher struct {
	msgChan chan *parser.Message
	sender  *sender.Manager
	workers int
	wg      sync.WaitGroup
}

func New(bufferSize, workers int, sender *sender.Manager) *Dispatcher {
	if workers <= 0 {
		workers = 2
	}
	if bufferSize <= 0 {
		bufferSize = 50
	}

	return &Dispatcher{
		msgChan: make(chan *parser.Message, bufferSize),
		sender:  sender,
		workers: workers,
	}
}

func (d *Dispatcher) Start(ctx context.Context) {
	for i := 0; i < d.workers; i++ {
		d.wg.Add(1)
		go d.worker(ctx)
	}
}

func (d *Dispatcher) Stop() {
	close(d.msgChan)
	d.wg.Wait()
}

func (d *Dispatcher) Dispatch(msg *parser.Message) {
	select {
	case d.msgChan <- msg:
		logger.Debug("Message dispatched",
			zap.String("platform", string(msg.Platform)),
			zap.String("content", msg.Content))
	default:
		logger.Error("Message channel is full, message dropped",
			zap.String("platform", string(msg.Platform)),
			zap.String("content", msg.Content))
	}
}

func (d *Dispatcher) worker(ctx context.Context) {
	defer d.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-d.msgChan:
			if !ok {
				return
			}
			if err := d.sender.Send(ctx, msg); err != nil {
				logger.Error("Failed to send message",
					zap.String("platform", string(msg.Platform)),
					zap.Error(err))
				continue
			}
			logger.Info("Message sent successfully",
				zap.String("platform", string(msg.Platform)))
		}
	}
}
