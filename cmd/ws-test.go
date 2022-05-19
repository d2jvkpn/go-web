package cmd

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/gorilla/websocket"
)

func NewWsTest() (command *cobra.Command) {
	var (
		addr string
		fSet *pflag.FlagSet
	)

	command = &cobra.Command{
		Use:   "ws-test",
		Short: "ws test",
		Long:  "websocket test by enter messages",

		Run: func(cmd *cobra.Command, args []string) {
			var (
				err  error
				link *url.URL
				conn *websocket.Conn
			)

			if !strings.HasPrefix(addr, "ws") {
				addr = "ws://" + addr
			}
			// link := url.URL{Scheme: "ws", Host: *addr, Path: "/ws/talk"}
			if link, err = url.Parse(addr); err != nil {
				log.Fatal("parse url:", err)
			}

			log.Printf("Connecting to %s\n", link.String())

			if conn, _, err = websocket.DefaultDialer.Dial(link.String(), nil); err != nil {
				log.Fatal("dial:", err)
			}
			defer conn.Close()

			runWs(conn)
		},
	}

	fSet = command.Flags()
	fSet.StringVar(
		&addr, "addr",
		"ws://localhost:8080/ws/v1/hello?name=Rover",
		"websocket service address",
	)

	return command
}

func runWs(conn *websocket.Conn) {
	var (
		pingId int
		mutex  *sync.Mutex
		ticker *time.Ticker
		done   chan struct{}
	)

	ticker = time.NewTicker(30 * time.Second)
	mutex = new(sync.Mutex) // avoid panic: concurrent write to websocket connection
	done = make(chan struct{})
	// send ping to server may be not necessary
	go func() {
	loop:
		for {
			select {
			case <-done:
				break loop
			case <-ticker.C:
				pingId++
				data := []byte(strconv.Itoa(pingId))
				log.Printf("~~> send ping: %q\n", data)
				mutex.Lock()
				_ = conn.WriteMessage(websocket.PingMessage, []byte(data))
				mutex.Unlock()
			}
		}
	}()

	// overwrite default handler(when receive a ping)
	conn.SetPingHandler(func(data string) (err error) {
		log.Printf("<~~ recv ping: %q, response pong\n", data)
		mutex.Lock()
		_ = conn.WriteMessage(websocket.PongMessage, []byte(data))
		mutex.Unlock()
		return nil
	})

	// overwrite default handler(when receive a pong after send a ping)
	conn.SetPongHandler(func(data string) (err error) {
		log.Printf("<~~ recv pong: %q\n", data)
		return nil
	})

	HandleMessage(conn, mutex, done)
	ticker.Stop()
}

func HandleMessage(conn *websocket.Conn, mutex *sync.Mutex, done chan struct{}) {
	var (
		typ int
		bts []byte
		err error
	)

	go func() {
		fmt.Println("Enter message and send to the server...")
	loop1:
		for {
			select {
			case <-done:
				break loop1
			default:
			}
			var msg string
			fmt.Scanf("%s", &msg)
			if msg == "" {
				close(done)
				break loop1
			}

			mutex.Lock()
			conn.WriteMessage(websocket.TextMessage, bytes.TrimSpace([]byte(msg)))
			mutex.Unlock()
		}
	}()

loop2:
	for {
		select {
		case <-done:
			break loop2
		default:
		}
		if typ, bts, err = conn.ReadMessage(); err != nil {
			log.Printf("<-- !!! ReadMessage error: %[1]T, %[1]v\n", err)
			break
		}
		log.Printf("<-- ReadMessage: type=%d\n\t%s\n", typ, bytes.TrimSpace(bts))
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
