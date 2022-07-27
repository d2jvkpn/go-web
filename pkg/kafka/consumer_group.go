package kafka

import (
	"context"
	"errors"
	"sync"

	"github.com/Shopify/sarama"
)

type Handler struct {
	Logger LogIntf

	group sarama.ConsumerGroup

	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup
}

func NewHandler(ctx context.Context, group sarama.ConsumerGroup) (handler *Handler) {
	handler = new(Handler)
	handler.group, handler.Logger = group, NewLogger("Kafka::ConsumerGroup")
	handler.ctx, handler.cancel = context.WithCancel(ctx)
	handler.wg = new(sync.WaitGroup)

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
				handler.Logger.Warn("<== Handler.Consume end")
			}

			if handler.ctx.Err() != nil {
				return
			}
		}
	}()
}

func (handler *Handler) Close() {
	handler.cancel()
	handler.wg.Wait()
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
				handler.Logger.Warn("~~~ Handle Setup done")
				return
			}
		}
	}()

	return nil
}

func (handler *Handler) Cleanup(sess sarama.ConsumerGroupSession) (err error) {
	// TODO
	handler.Logger.Info(">>> Handler Cleanup")
	return nil
}

func (handler *Handler) ConsumeClaim(sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) (err error) {

	handler.wg.Add(1)
	defer handler.wg.Done()

	tmpl := "<-- msg.Timestamp=%+v, msg.Topic=%v, msg.Partition=%v, msg.Offset=%v\n" +
		"    key: %q, value: %q\n"

LOOP:
	for {
		select {
		case msg := <-claim.Messages():
			if msg == nil {
				break LOOP
			}
			handler.Logger.Info(
				tmpl, msg.Timestamp, msg.Topic, msg.Partition, msg.Offset,
				msg.Key, msg.Value,
			)

			// TODO: process(msg)
			// sess.MarkOffset(msg.Topic, msg.Partition, msg.Offset, "some-metadata")
			sess.MarkMessage(msg, "consumed-by-d2jvkpn")
		case <-handler.ctx.Done():
			handler.Logger.Warn("!!! ConsumeClaim canceled")
			break LOOP
		}
	}

	return nil
}
