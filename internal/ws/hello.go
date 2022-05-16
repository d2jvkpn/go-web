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

func hello(ctx *gin.Context) {
	var (
		// client information
		clientId    uint64    = atomic.AddUint64(&_ClientId, 1)
		name        string    = ctx.DefaultQuery("name", "World")
		ip          string    = ctx.ClientIP()
		connectTime time.Time = time.Now()
		pongTime    time.Time
	)
	_, _ = connectTime, pongTime

	id := fmt.Sprintf("name=%s, ip=%s, clientId=%d", name, ip, clientId)
	mel, once := New(), new(sync.Once)
	// log.Printf("%+v\n", mel.Config)
	mel.Config.PingPeriod = 10 * time.Second

	mel.HandleConnect(func(sess *Session) {
		log.Printf(">>> hello new ws connection: %q\n", id)
	})

	mel.HandleDisconnect(func(sess *Session) {
		log.Printf("<<< hello ws disconnected: %q\n", id)
	})

	mel.HandleError(func(sess *Session, err error) {
		log.Printf("!!! hello ws error: %q, error=%q\n", id, err)
	})

	mel.HandlePong(func(sess *Session) {
		pongTime = time.Now()
		log.Printf("<~~ %q recv pong\n", id)
	})

	mel.HandleMessage(func(sess *Session, msg []byte) {
		once.Do(func() {
			data := fmt.Sprintf(`{"type":"clientId","clientId":%d}`, clientId)
			_ = sess.Write([]byte(data))
		})

		// m.Broadcast(msg)
		log.Printf("<-- %q recv: %q\n", id, msg)
		send := fmt.Sprintf("%s, nice to meet you!", name)
		log.Printf("--> %q send: %q\n", id, send)
		_ = sess.Write([]byte(send))
	})

	// _ = mel.HandleRequest(ctx.Writer, ctx.Request)
	_ = mel.HandleRequestWithKeys(ctx.Writer, ctx.Request, map[string]interface{}{
		"name": name, "ip": ctx.ClientIP(),
	})
}
