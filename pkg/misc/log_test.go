package misc

import (
	"io"
	"log"
	"testing"
)

func TestLogWriter(t *testing.T) {
	lw, err := NewLogWriter("logs/test", io.Discard)
	if err != nil {
		t.Fatal(err)
	}

	lw.Register()

	log.Println("hello, world!")
	log.Println("INFO")

	lw.Close()
	log.Println("XXXX")
}

func TestRegisterDefaultLogFmt(t *testing.T) {
	RegisterDefaultLogFmt()
	log.Printf("Hello, %s", "Rover")
	log.Printf("Hello, %s\n", "d2jvkpn")
}
