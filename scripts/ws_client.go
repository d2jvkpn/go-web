package main

import (
	"bytes"
	"flag"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	var (
		pingId int
		err    error
		link   *url.URL
		ticker *time.Ticker
		conn   *websocket.Conn
	)

	addr := flag.String(
		"addr",
		"ws://localhost:8080/ws/v1/hello?name=Rover",
		"websocket service address",
	)

	flag.Parse()

	// link := url.URL{Scheme: "ws", Host: *addr, Path: "/ws/talk"}
	if link, err = url.Parse(*addr); err != nil {
		log.Fatal("parse url:", err)
	}
	log.Printf("Connecting to %s\n", link.String())

	if conn, _, err = websocket.DefaultDialer.Dial(link.String(), nil); err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	ticker = time.NewTicker(30 * time.Second)
	// send ping to server may be not necessary
	go func() {
		for {
			select {
			case <-ticker.C:
				pingId++
				data := []byte(strconv.Itoa(pingId))
				log.Printf("~~> send ping: %q\n", data)
				_ = conn.WriteMessage(websocket.PingMessage, []byte(data))
			}
		}
	}()

	// overwrite default handler(when receive a ping)
	conn.SetPingHandler(func(data string) (err error) {
		log.Printf("<~~ recv ping: %q, response pong\n", data)
		_ = conn.WriteMessage(websocket.PongMessage, []byte(data))
		return nil
	})

	// overwrite default handler(when receive a pong after send a ping)
	conn.SetPongHandler(func(data string) (err error) {
		log.Printf("<~~ recv pong: %q\n", data)
		return nil
	})

	HandleMessage(conn)
	ticker.Stop()
}

func HandleMessage(conn *websocket.Conn) {
	var (
		typ int
		bts []byte
		err error
	)

	msg := "Hello, I'm a novice."
	log.Printf("--> WriteMessage: %q\n", msg)
	conn.WriteMessage(websocket.TextMessage, []byte(msg))

	for {
		if typ, bts, err = conn.ReadMessage(); err != nil {
			log.Printf("<-- !!! ReadMessage error: %[1]T, %[1]v\n", err)
			break
		}
		log.Printf("<-- ReadMessage: type=%d, msg=%q\n", typ, bytes.TrimSpace(bts))
		if typ == websocket.CloseMessage {
			break
		}
	}
}

// https://pkg.go.dev/github.com/gorilla/websocket#pkg-types
//	TextMessage = 1
//	BinaryMessage = 2
//	CloseMessage = 8
//	PingMessage = 9
//	PongMessage = 10