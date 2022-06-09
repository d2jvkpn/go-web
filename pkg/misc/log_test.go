package misc

import (
	"log"
	"testing"
)

func TestLogWriter(t *testing.T) {
	lw, err := NewLogWriter("logs/test", true)
	if err != nil {
		t.Fatal(err)
	}

	lw.Register()

	log.Println("hello, world!")
	log.Println("INFO")

	lw.Close()
	log.Println("XXXX")
}
