package teerotate

import (
	"io"
)

// tmpLogWriter is short lived & immutable
// it corresponds to a log file
// and is rotated out
type tmpLogWriter struct {
	w                    io.WriteCloser
	rotateConditionCheck rotateConditionCheck
}

func newTmpLogger(w io.WriteCloser, rotateConditionCheck rotateConditionCheck) *tmpLogWriter {
	return &tmpLogWriter{
		w:                    w,
		rotateConditionCheck: rotateConditionCheck,
	}
}

func (l tmpLogWriter) isDone() bool {
	return l.rotateConditionCheck()
}

func (l tmpLogWriter) Write(b []byte) (done bool) {
	l.w.Write(b)
	return l.isDone()
}

func (l tmpLogWriter) Close() {
	l.w.Close()
}
