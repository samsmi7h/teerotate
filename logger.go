package teerotate

import (
	"io"
	"time"
)

// tmpLogWriter is short lived & immutable
// it corresponds to a log file
// and is rotated out
type tmpLogWriter struct {
	w      io.WriteCloser
	ticker <-chan time.Time
}

func newTmpLogger(w io.WriteCloser, t <-chan time.Time) *tmpLogWriter {
	return &tmpLogWriter{
		w:      w,
		ticker: t,
	}
}

func (l tmpLogWriter) isDone() bool {
	select {
	case <-l.ticker:
		return true
	default:
		// continue
	}

	return false
}

func (l tmpLogWriter) Write(b []byte) (done bool) {
	l.w.Write(b)
	return l.isDone()
}

func (l tmpLogWriter) Close() {
	l.w.Close()
}
