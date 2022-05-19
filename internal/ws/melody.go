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
	"gopkg.in/olahol/melody.v1"
)

var (
	_ClientId uint64
	_MelHello *melody.Melody
)

func init() {
	_MelHello = melody.New()
	// log.Printf("%+v\n", _MelHello.Config)
	_MelHello.Config.PingPeriod = 10 * time.Second
	_MelHello.Upgrader.CheckOrigin = func(req *http.Request) bool { return true }

	_MelHello.HandleConnect(func(sess *melody.Session) {
		var client *Client
		if client = getSessionClient(sess); client == nil {
			return
		}

		log.Printf(">>> hello new ws connection: %q, addr=%s\n", client, client.Address)
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

	_MelHello.HandleDisconnect(func(sess *melody.Session) {
		var client *Client
		if client = getSessionClient(sess); client == nil {
			return
		}

		log.Printf("<<< hello ws disconnected: %q\n", client)
	})

	_MelHello.HandleError(func(sess *melody.Session, err error) {
		var client *Client
		if client = getSessionClient(sess); client == nil {
			return
		}

		log.Printf("!!! hello ws error: %q, error=%q\n", client, err)
	})

	_MelHello.HandlePong(func(sess *melody.Session) {
		var client *Client
		if client = getSessionClient(sess); client == nil {
			return
		}

		client.PongTime = time.Now()
		log.Printf("<~~ %q recv pong\n", client)
	})

	_MelHello.HandleMessage(func(sess *melody.Session, msg []byte) {
		var client *Client
		if client = getSessionClient(sess); client == nil {
			return
		}

		if len(msg) == 0 {
			return
		}

		var (
			data []byte
			err  error
		)

		log.Printf("<-- %q recv: %q\n", client, msg)
		if data, err = client.HandleMsg(msg); err != nil {
			log.Printf("%s HandleMsg error and close connection: %v\n", client, err)
			_ = sess.CloseWithMsg([]byte(err.Error()))
			return
		}

		if len(data) > 0 {
			_ = sess.Write(data) // TODO: verify concurrent writing safe
		}
		// _MelHello.Broadcast(msg)
	})
}

type Client struct {
	Id            string    `json:"id"`
	Address       string    `json:"address"`
	Name          string    `json:"name"`
	ConnectTime   time.Time `json:"connectTime"`
	PongTime      time.Time `json:"pongTime"`
	melo          *melody.Melody
	*sync.RWMutex `json:"-"`
}

func NewClient(addr, name string, melo *melody.Melody) *Client {
	return &Client{
		Id:          fmt.Sprintf("%04d", atomic.AddUint64(&_ClientId, 1)),
		Address:     addr,
		Name:        name,
		ConnectTime: time.Now(),
		melo:        melo,
		RWMutex:     new(sync.RWMutex),
	}
}

func (client Client) String() string {
	return fmt.Sprintf("name=%s, id=%s", client.Name, client.Id)
}

func (client *Client) HandleMsg(msg []byte) (data []byte, err error) {
	// more processes...
	return msg, nil
}

func getSessionClient(sess *melody.Session) (client *Client) {
	var exists bool

	if sess.IsClosed() {
		return nil
	}

	if client, exists = sess.Keys["client"].(*Client); !exists {
		return nil
	}

	return client
}
