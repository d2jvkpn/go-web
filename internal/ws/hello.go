package ws

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	melody "gopkg.in/olahol/melody.v1"
)

func hello(ctx *gin.Context) {
	mel := melody.New()
	name := ctx.DefaultQuery("name", "world")

	mel.HandleMessage(func(sess *melody.Session, msg []byte) {
		// m.Broadcast(msg)
		log.Printf("<-- recv: %q\n", msg)
		send := fmt.Sprintf("%s, nice to meet you!", name)
		log.Printf("<-- send: %q\n", send)
		sess.Write([]byte(send))
	})

	mel.HandleRequest(ctx.Writer, ctx.Request)
}
