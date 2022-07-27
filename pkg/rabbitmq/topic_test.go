package rabbitmq

import (
	"fmt"
	"testing"
)

func TestTopicQueue(t *testing.T) {
	q, err := NewTopicQueue(testUri, "articles", "collect", "xxxxxx")

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(q.Q)
}

func TestTopicSender_t1(t *testing.T) {
	sender, err := NewTopicSender(testUri, "articles", "collect")
	if err != nil {
		t.Fatal(err)
	}

	if err := sender.Emit([]byte("hello")); err != nil {
		t.Fatal(err)
	}
}

func TestTopicReceiver_t1(t *testing.T) {
	rec, err := NewTopicReceiver(testUri, "articles", "collect", "collect_link", "")
	if err != nil {
		t.Fatal(err)
	}

	d := <-rec.D
	fmt.Println(string(d.Body))
}
