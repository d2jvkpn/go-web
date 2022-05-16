package ws

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	. "gopkg.in/olahol/melody.v1"
)

var (
	_ClientId uint64
)

type Client struct {
	Id          uint64
	Ip          string
	Name        string
	ConnectTime time.Time
	PongTime    time.Time
}

func NewClient(id uint64, ip, name string) *Client {
	return &Client{
		Id:          id,
		Ip:          ip,
		Name:        name,
		ConnectTime: time.Now(),
	}
}

func (client Client) String() string {
	return fmt.Sprintf("name=%s, id=%d", client.Name, client.Id)
}

func hello(ctx *gin.Context) {
	client := NewClient(
		atomic.AddUint64(&_ClientId, 1),
		ctx.ClientIP(),
		ctx.DefaultQuery("name", "World"),
	)

	mel, once := New(), new(sync.Once)
	// log.Printf("%+v\n", mel.Config)
	mel.Config.PingPeriod = 10 * time.Second

	mel.HandleConnect(func(sess *Session) {
		log.Printf(">>> hello new ws connection: %q\n", client)
	})

	mel.HandleDisconnect(func(sess *Session) {
		log.Printf("<<< hello ws disconnected: %q\n", client)
	})

	mel.HandleError(func(sess *Session, err error) {
		log.Printf("!!! hello ws error: %q, error=%q\n", client, err)
	})

	mel.HandlePong(func(sess *Session) {
		client.PongTime = time.Now()
		log.Printf("<~~ %q recv pong\n", client)
	})

	mel.HandleMessage(func(sess *Session, msg []byte) {
		once.Do(func() {
			data := fmt.Sprintf(`{"type":"clientId","clientId":%d}`, client.Id)
			_ = sess.Write([]byte(data))
		})

		// m.Broadcast(msg)
		send := fmt.Sprintf("%s, nice to meet you!", client.Name)
		log.Printf("<-- %q recv: %q, send: %q\n", client, msg, send)
		_ = sess.Write([]byte(send))
	})

	// _ = mel.HandleRequest(ctx.Writer, ctx.Request)
	_ = mel.HandleRequestWithKeys(
		ctx.Writer, ctx.Request,
		map[string]interface{}{"client": client},
	)
}
