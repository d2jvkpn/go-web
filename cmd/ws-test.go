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

	"github.com/d2jvkpn/goapp/pkg/misc"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
			// defer conn.Close()

			client := NewWsClient(conn, 30*time.Second, true)
			client.HandleMessage()
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

type WsClient struct {
	conn    *websocket.Conn
	mutex   *sync.Mutex
	done    chan struct{}
	pingDur time.Duration
	ticker  *time.Ticker
	jsonMsg bool
}

func NewWsClient(conn *websocket.Conn, pingDur time.Duration, jsonMsg bool) (client WsClient) {
	client = WsClient{
		conn:    conn,
		mutex:   new(sync.Mutex),
		done:    make(chan struct{}),
		pingDur: pingDur,
		jsonMsg: jsonMsg,
	}

	// overwrite default handler(when receive a ping)
	conn.SetPingHandler(func(data string) (err error) {
		log.Printf("<~~ recv ping: %q, response pong\n", data)
		client.mutex.Lock()
		_ = conn.WriteMessage(websocket.PongMessage, []byte(data))
		client.mutex.Unlock()
		return nil
	})

	// overwrite default handler(when receive a pong after send a ping)
	conn.SetPongHandler(func(data string) (err error) {
		log.Printf("<~~ recv pong: %q\n", data)
		return nil
	})

	if pingDur <= 0 {
		return client
	}
	client.ticker = time.NewTicker(pingDur)

	go func() {
		var pingId int
	loop:
		for {
			select {
			case <-client.done:
				break loop
			case <-client.ticker.C:
				pingId++
				data := []byte(strconv.Itoa(pingId))
				log.Printf("~~> send ping: %q\n", data)
				client.mutex.Lock()
				_ = client.conn.WriteMessage(websocket.PingMessage, []byte(data))
				client.mutex.Unlock()
			}
		}
	}()

	return client
}

func (client WsClient) Close() {
	client.ticker.Stop()
	close(client.done)
	client.conn.Close()
}

func (client WsClient) HandleMessage() {
	var (
		typ int
		bts []byte
		err error
	)

	go func() {
		fmt.Println(">>> Enter message and send to the server...")
		for {
			var (
				bts []byte
				msg string
			)
			fmt.Scanf("%s", &msg)
			msg = strings.TrimSpace(msg)
			if msg == "\\q" {
				log.Println("!!! Exit Client")
				client.Close()
				break
			}
			if bts = []byte(msg); len(bts) == 0 {
				continue
			}

			if client.jsonMsg && misc.CheckJson(bts) != nil {
				log.Printf("!!! invalid json: %q\n", msg)
				continue
			}
			log.Printf("<<< %q\n", msg)

			client.mutex.Lock()
			client.conn.WriteMessage(websocket.TextMessage, bts)
			client.mutex.Unlock()
		}
	}()

	for {
		if typ, bts, err = client.conn.ReadMessage(); err != nil {
			log.Printf("!!! ReadMessage error: %[1]T, %[1]v\n", err)
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
