package teerotate

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type closerWrapper struct {
	io.ReadWriter
	closed bool
}

func (c *closerWrapper) Close() error {
	c.closed = true
	return nil
}

func closer(rw io.ReadWriter) *closerWrapper {
	return &closerWrapper{ReadWriter: rw}
}

func TestRotate(t *testing.T) {
	buffers := []*bytes.Buffer{}
	done := make(chan time.Time)

	l := newRotatingLogger(
		func() io.WriteCloser {
			buf := bytes.NewBuffer([]byte{})
			buffers = append(buffers, buf)
			return closer(buf)
		},
		func() <-chan time.Time {
			return done
		},
	)

	l.Print("hello")
	done <- time.Now()
	l.Print("hallo")
	l.Print("world")
	l.Close()

	assert.Equal(t, len(buffers), 2)

	firstBuf, _ := io.ReadAll(buffers[0])
	secondBuf, _ := io.ReadAll(buffers[1])
	b := append(firstBuf, secondBuf...)
	assert.Equal(t, "hellohalloworld", string(b))
}
