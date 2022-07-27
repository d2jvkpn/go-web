package x

import (
	"fmt"

	"github.com/streadway/amqp"
)

type TopicReceiver struct {
	TopicQueue
	D <-chan amqp.Delivery
}

func (queue *TopicQueue) NewReceiver(customer string) (rece *TopicReceiver, err error) {
	rece = &TopicReceiver{TopicQueue: *queue}

	rece.D, err = rece.ch.Consume(
		rece.Q,   // queue
		customer, // consumer
		false,    // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)

	return
}

func (rec *TopicReceiver) Read() (bts []byte, err error) {
	select {
	case d := <-rec.D:
		bts = d.Body
	default:
		err = fmt.Errorf("no new message")
	}

	return
}

// create Receiver
func NewTopicReceiver(uri, exchange, routingKey, queueName, customer string) (
	rec *TopicReceiver, err error) {
	var queue *TopicQueue

	if queue, err = NewTopicQueue(uri, exchange, routingKey, queueName); err != nil {
		return nil, err
	}

	rec = &TopicReceiver{TopicQueue: *queue}
	rec.D, err = rec.ch.Consume(
		rec.Q,    // queue
		customer, // consumer
		false,    // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)
	if err != nil {
		return nil, err
	}

	return
}
