package rabbitmq

import (
	// "fmt"

	"github.com/streadway/amqp"
)

type TopicSender struct {
	E, R string
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewTopicSender(uri, exchange, routingKey string) (sender *TopicSender, err error) {

	// addr := fmt.Sprintf("amqp://%s:%s@%s:%s%s", user, passwd, host, port, vhost)
	sender = new(TopicSender)
	if sender.conn, err = amqp.Dial(uri); err != nil {
		return nil, err
	}

	if sender.ch, err = sender.conn.Channel(); err != nil {
		return nil, err
	}

	sender.E, sender.R = exchange, routingKey

	err = sender.ch.ExchangeDeclare(
		sender.E, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)

	return
}

func (sender *TopicSender) Emit(bts []byte) (err error) {
	err = sender.ch.Publish(
		sender.E, // exchange
		sender.R, // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        bts,
		})

	return
}

func (sender *TopicSender) Close() {
	// if sender.ch != nil {
	//	sender.ch.Close()
	// }

	if sender.conn != nil {
		sender.conn.Close()
	}
}
