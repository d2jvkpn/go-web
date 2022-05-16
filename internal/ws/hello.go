package ws

import (
	"encoding/json"
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
	_Mel      *Melody = New()
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

	once := new(sync.Once)
	// log.Printf("%+v\n", _Mel.Config)
	_Mel.Config.PingPeriod = 10 * time.Second

	_Mel.HandleConnect(func(sess *Session) {
		log.Printf(">>> hello new ws connection: %q, ip=%s\n", client, client.Ip)
	})

	_Mel.HandleDisconnect(func(sess *Session) {
		log.Printf("<<< hello ws disconnected: %q\n", client)
	})

	_Mel.HandleError(func(sess *Session, err error) {
		log.Printf("!!! hello ws error: %q, error=%q\n", client, err)
	})

	_Mel.HandlePong(func(sess *Session) {
		client.PongTime = time.Now()
		log.Printf("<~~ %q recv pong\n", client)
	})

	_Mel.HandleMessage(func(sess *Session, msg []byte) {
		log.Printf("<-- %q recv: %q\n", client, msg)

		once.Do(func() {
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

	// _ = _Mel.HandleRequest(ctx.Writer, ctx.Request)
	_ = _Mel.HandleRequestWithKeys(
		ctx.Writer, ctx.Request,
		map[string]interface{}{"client": client},
	)
}
