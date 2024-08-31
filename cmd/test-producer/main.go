package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/samsmi7h/teerotate"
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	l := teerotate.NewRotatingFileLogger("/tmp", time.Hour)

	for {
		select {
		case <-sigChan:
			fmt.Println("got signal")
			l.Close()
			fmt.Println("finished")
			return
		default:
			l.Print(time.Now().String() + "\n")
		}
	}
}
