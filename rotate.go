package teerotate

import (
	"fmt"
	"os"
	"time"
)

type RotatingLogger struct {
	newOuput  makeNewOutput
	newTicker makeNewTicker

	printCh chan []byte
	doneCh  chan struct{}
	hooks   hooks

	currentLogger *tmpLogWriter
}

// Print is the method exposed for printing a log.
// It hands off handling of the message.
// But will block if the queue is full.
// Then prints to stdout
func (r *RotatingLogger) Print(msg string, args ...interface{}) {
	s := fmt.Sprintf(msg, args...)
	r.printCh <- []byte(s)
	fmt.Fprint(os.Stdout, s)
}

// worker that writes logs
func (r *RotatingLogger) startPrinter() {
	for b := range r.printCh {
		done := r.currentLogger.Write(b)
		if done {
			r.Rotate()
		}
	}

	fmt.Println("closing current logger...")
	r.currentLogger.Close()
	fmt.Println("current logger closed.")
	r.doneCh <- struct{}{}
}

func (r *RotatingLogger) Rotate() {
	w := r.newOuput()
	t := r.newTicker()
	newLogger := newTmpLogger(w, t)

	// swap in the new one
	oldLogger := r.currentLogger
	r.currentLogger = newLogger

	r.currentLogger = newLogger

	// is nil at start
	if oldLogger != nil {
		oldLogger.Close()

		if r.hooks.postRotation != nil {
			go r.hooks.postRotation()
		}
	}
}

func (r *RotatingLogger) Close() {
	close(r.printCh)
	<-r.doneCh

	if r.hooks.postRotation != nil {
		// TODO: test this
		r.hooks.postRotation()
	}
}

func NewRotatingFileLogger(dir string, lifespan time.Duration) *RotatingLogger {
	return newRotatingLogger(
		fileFactory(dir),
		tickerFactory(lifespan),
	)
}

// newRotatingLogger is lower level and testable
func newRotatingLogger(no makeNewOutput, nt makeNewTicker) *RotatingLogger {
	r := RotatingLogger{
		newOuput:  no,
		newTicker: nt,
		printCh:   make(chan []byte, 1000),
		doneCh:    make(chan struct{}, 1),
	}

	r.Rotate()
	go r.startPrinter()

	return &r
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
