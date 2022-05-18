package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	. "gopkg.in/olahol/melody.v1"
)

var (
	_ClientId uint64
	_MelHello *Melody = New()
)

type Client struct {
	Id          string
	Ip          string
	Name        string
	ConnectTime time.Time
	PongTime    time.Time
}

func NewClient(id, ip, name string) *Client {
	return &Client{
		Id:          id,
		Ip:          ip,
		Name:        name,
		ConnectTime: time.Now(),
	}
}

func (client Client) String() string {
	return fmt.Sprintf("name=%s, id=%s", client.Name, client.Id)
}

func hello(ctx *gin.Context) {
	var (
		once   *sync.Once
		client *Client
	)

	once = new(sync.Once)

	client = NewClient(
		fmt.Sprintf("%04d", atomic.AddUint64(&_ClientId, 1)),
		ctx.ClientIP(),
		ctx.DefaultQuery("name", "World"),
	)

	// log.Printf("%+v\n", _MelHello.Config)
	_MelHello.Config.PingPeriod = 10 * time.Second
	_MelHello.Upgrader.CheckOrigin = func(req *http.Request) bool { return true }

	_MelHello.HandleConnect(func(sess *Session) {
		log.Printf(">>> hello new ws connection: %q, ip=%s\n", client, client.Ip)
	})

	_MelHello.HandleDisconnect(func(sess *Session) {
		log.Printf("<<< hello ws disconnected: %q\n", client)
	})

	_MelHello.HandleError(func(sess *Session, err error) {
		log.Printf("!!! hello ws error: %q, error=%q\n", client, err)
	})

	_MelHello.HandlePong(func(sess *Session) {
		client.PongTime = time.Now()
		log.Printf("<~~ %q recv pong\n", client)
	})

	_MelHello.HandleMessage(func(sess *Session, msg []byte) {
		log.Printf("<-- %q recv: %q\n", client, msg)

		once.Do(func() {
			var client *Client
			if sess.Keys != nil {
				client, _ = sess.Keys["client"].(*Client)
			}
			log.Printf("~~~ client: %+v\n", client)

			data := map[string]interface{}{
				"code":    0,
				"type":    "client",
				"message": fmt.Sprintf("%s, nice to meet you!", client.Name),
				"data":    gin.H{"clientId": client.Id},
			}

			log.Printf("--> %q send client information\n", client)
			bts, _ := json.Marshal(data)
			_ = sess.Write(bts)
		})

		// m.Broadcast(msg)
	})

	// _ = _MelHello.HandleRequest(ctx.Writer, ctx.Request)
	_ = _MelHello.HandleRequestWithKeys(
		ctx.Writer, ctx.Request, map[string]interface{}{"client": client},
	)
}
