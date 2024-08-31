package teerotate

import (
	"fmt"
	"os"
	"time"
)

// TODO: post rotate hook e.g. S3 upload
// TODO: could even pass in your own logger

type RotatingLogger struct {
	newOuput  makeNewOutput
	newTicker makeNewTicker
	ch        chan []byte

	currentLogger *tmpLogWriter
}

// Print is the method exposed for printing a log.
// It hands off handling of the message.
// But will block if the queue is full.
// Then prints to stdout
func (r *RotatingLogger) Print(msg string, args ...interface{}) {
	s := fmt.Sprintf(msg, args...)
	r.ch <- []byte(s)
	fmt.Fprint(os.Stdout, s)
}

// worker that writes logs
func (r *RotatingLogger) startPrinter() {
	for b := range r.ch {
		done := r.currentLogger.Write(b)
		if done {
			r.Rotate()
		}
	}
}

func (r *RotatingLogger) Rotate() {
	w := r.newOuput()
	t := r.newTicker()
	newLogger := newTmpLogger(w, t)

	// swap in the new one
	oldLogger := r.currentLogger
	r.currentLogger = newLogger

	r.currentLogger = newLogger

	// nil at start
	if oldLogger != nil {
		oldLogger.Close()
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
		ch:        make(chan []byte, 1000),
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
