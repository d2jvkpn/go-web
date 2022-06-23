package internal

import (
	// "fmt"
	"log"
	"os"

	"github.com/d2jvkpn/go-web/pkg/misc"
)

func Run(config, addr string, release bool) {
	var (
		value int
		err   error
		ch    chan int
	)

	if err = Load(config, release); err != nil {
		log.Fatalln(err)
	}

	ch = make(chan int, 1)
	go misc.ListenOSSignal(ch) // goroutine1, send 0 when interrupted, send -1 otherwise

	go func() { // goroutine2
		if err = Serve(addr); err != nil {
			log.Println(err)
			ch <- 1
		} else {
			ch <- 2
		}

	}()

	value = <-ch
	switch value {
	case 1, 2: // goroutine2 exit
		Down()
		ch <- value // send to goroutine1
		<-ch        // goroutine1 exit: -1
	case 0: // both goroutine1 exit by interrupted
		Down()
		<-ch // goroutine2 exit: 1
	case -1: // both goroutine1 and goroutine2 exited
		Down()
	}

	if value != 0 {
		os.Exit(1)
	}
}
