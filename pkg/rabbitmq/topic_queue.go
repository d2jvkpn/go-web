package x

import (
	"fmt"

	"github.com/streadway/amqp"
)

type TopicQueue struct {
	E, R, Q string // exchange, routingKey, queueName
	conn    *amqp.Connection
	ch      *amqp.Channel
}

// create Receiver
// uri: fmt.Sprintf("amqp://%s:%s@%s:%s%s", user, passwd, host, port, vhost)
func NewTopicQueue(uri, exchange, routingKey, queueName string) (queue *TopicQueue, err error) {
	var q amqp.Queue

	if queueName == "" {
		return nil, fmt.Errorf("queueName is empty")
	}

	queue = &TopicQueue{
		E: exchange,
		R: routingKey,
		Q: queueName,
	}

	if queue.conn, err = amqp.Dial(uri); err != nil {
		return nil, err
	}

	if queue.ch, err = queue.conn.Channel(); err != nil {
		return nil, err
	}

	err = queue.ch.ExchangeDeclare(
		queue.E, // name
		"topic", // type
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return nil, err
	}

	q, err = queue.ch.QueueDeclare(
		queue.Q, // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return nil, err
	}

	err = queue.ch.QueueBind(
		q.Name,  // queue name
		queue.R, // routing key
		queue.E, // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return
}

func (queue *TopicQueue) Close() {
	//if queue.ch != nil {
	//	queue.ch.Close()
	//	}

	if queue.conn != nil {
		queue.conn.Close()
	}
}
