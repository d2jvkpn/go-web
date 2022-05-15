package ws

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	. "gopkg.in/olahol/melody.v1"
)

func hello(ctx *gin.Context) {
	mel := New()
	name := ctx.DefaultQuery("name", "World")
	id := fmt.Sprintf("name=%s, ip=%s", name, ctx.ClientIP())

	mel.HandleConnect(func(sess *Session) {
		log.Printf(">>> hello new ws connection: %q\n", id)
	})

	mel.HandleDisconnect(func(sess *Session) {
		log.Printf("<<< hello ws disconnected: %q\n", id)
	})

	mel.HandleError(func(sess *Session, err error) {
		log.Printf("!!! hello ws error: %q, error=%q\n", id, err)
	})

	mel.HandleMessage(func(sess *Session, msg []byte) {
		// m.Broadcast(msg)
		log.Printf("<-- %q recv: %q\n", id, msg)
		send := fmt.Sprintf("%s, nice to meet you!", name)
		log.Printf("--> %q send: %q\n", id, send)
		sess.Write([]byte(send))
	})

	mel.HandleRequest(ctx.Writer, ctx.Request)
}
