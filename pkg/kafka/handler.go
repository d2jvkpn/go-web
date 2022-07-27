package kafka

import (
	"context"
	"errors"
	"sync"

	"github.com/Shopify/sarama"
)

type Handler struct {
	Logger LogIntf
	// process a message, MarkMessage(msg, metadata) as consumed if return string isn't empty
	process Process
	group   sarama.ConsumerGroup

	wg     *sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

type Process func(msg *sarama.ConsumerMessage) (metadata string, err error)

func NewHandler(ctx context.Context, group sarama.ConsumerGroup, process Process) (
	handler *Handler) {

	handler = &Handler{
		Logger:  NewLogger(),
		process: process,
		group:   group,
		wg:      new(sync.WaitGroup),
	}

	handler.ctx, handler.cancel = context.WithCancel(ctx)

	return handler
}

func (handler *Handler) Consume(topics ...string) {
	go func() {
		var err error
		handler.Logger.Info("==> Handler.Consume start")
		for {
			err = handler.group.Consume(handler.ctx, topics, handler)
			if err != nil {
				if errors.Is(sarama.ErrClosedConsumerGroup, err) {
					handler.Logger.Warn("!!! Handler.Consume closed:", err)
					return
				}
				handler.Logger.Error("!!! Handler.Consume error:", err)
			} else {
				handler.Logger.Warn("<== Handler.Consume end") // occurs when reset offset
			}

			if err = handler.ctx.Err(); err != nil {
				handler.Logger.Error("!!! Handler.Consume ctx.Err(): %v", err)
				return
			}
		}
	}()
}

func (handler *Handler) Close() error {
	handler.cancel()
	handler.wg.Wait()

	return handler.group.Close()
}

func (handler *Handler) Setup(sess sarama.ConsumerGroupSession) (err error) {
	handler.Logger.Info(">>> Handler.Setup Start")

	go func() {
		var err error
		for {
			select {
			case err = <-handler.group.Errors():
				if errors.Is(sarama.ErrClosedConsumerGroup, err) {
					handler.Logger.Warn("!!! Handle.Setup closed")
					return
				}
				handler.Logger.Error("!!! Handle.Setup error:", err)
			case <-handler.ctx.Done():
				handler.Logger.Info("~~~ Handle.Setup done")
				return
			}
		}
	}()

	return nil
}

func (handler *Handler) Cleanup(sess sarama.ConsumerGroupSession) (err error) {
	// TODO
	handler.Logger.Info(">>> Handler.Cleanup start")
	return nil
}

func (handler *Handler) ConsumeClaim(sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {

	handler.wg.Add(1)
	defer handler.wg.Done()

LOOP:
	for {
		select {
		case msg := <-claim.Messages():
			if msg == nil {
				break LOOP
			}

			metadata, err := handler.process(msg)
			if err != nil {
				handler.Logger.Error("!!! Handler.ConsumeClaim process: %v", err)
			}
			if metadata != "" {
				// sess.MarkOffset(msg.Topic, msg.Partition, msg.Offset, "some-metadata")
				sess.MarkMessage(msg, metadata)
			}
		case <-handler.ctx.Done():
			handler.Logger.Warn("!!! Handler.ConsumeClaim canceled")
			break LOOP
		}
	}

	return nil
}
