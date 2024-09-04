package teerotate

import (
	"bytes"
	"io"
	"strings"
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
		bytes.NewBuffer([]byte{}),
	)

	hookCh := make(chan struct{}, 2)
	l.WithPostRotationHook(func() {
		hookCh <- struct{}{}
	})

	l.Print("hello\n")
	done <- time.Now()

	// hook should fire at rotate
	assert.True(t, wasChanWrittenTo(hookCh))

	l.Print("hallo\n")
	l.Print("world\n")
	l.Close()

	// hook should fire at end
	assert.True(t, wasChanWrittenTo(hookCh))

	assert.Equal(t, len(buffers), 2)

	firstBuf, _ := io.ReadAll(buffers[0])
	secondBuf, _ := io.ReadAll(buffers[1])
	b := append(firstBuf, secondBuf...)
	assert.Equal(t, "hello\nhallo\nworld\n", string(b))
}

func wasChanWrittenTo[T any](ch chan T) bool {
	select {
	case <-ch:
		return true
	default:
		return false
	}
}

type writerWithDelay struct {
	io.WriteCloser
	delay time.Duration
}

func (w writerWithDelay) Write(p []byte) (n int, err error) {
	time.Sleep(w.delay)
	return w.WriteCloser.Write(p)
}

func newWriterWithDelay(w io.WriteCloser, t time.Duration) writerWithDelay {
	return writerWithDelay{
		WriteCloser: w,
		delay:       t,
	}
}

func TestCloseGetsAllLogs(t *testing.T) {
	type testCase struct {
		callsClose          bool
		assertNumberOfLines func(t *testing.T, expected, actual int)
		assertLastLine      func(t *testing.T, expected, actual string)
	}

	tcs := []testCase{
		{
			// WHEN WE DONT CALL CLOSE
			// THE BACKLOG DOESNT GET CLEARED
			// MEANING WERE MISSING LOGS FROM THE FINAL OUTPUT
			callsClose: false,
			assertNumberOfLines: func(t *testing.T, expected, actual int) {
				assert.Less(t, actual, expected)
			},
			assertLastLine: func(t *testing.T, expected, actual string) {
				assert.NotEqual(t, expected, actual)
			},
		},

		{
			// WHEN WE DO CALL CLOSE
			// THE BACKLOG GET DRAINED
			// MEANING WE DONT MISS ANYTHING
			callsClose: true,
			assertNumberOfLines: func(t *testing.T, expected, actual int) {
				assert.Equal(t, expected, actual)
			},
			assertLastLine: func(t *testing.T, expected, actual string) {
				assert.Equal(t, expected, actual)
			},
		},
	}

	const numberOfLogs = 1000

	for _, tc := range tcs {
		output := bytes.NewBuffer([]byte{})
		l := newRotatingLogger(
			func() io.WriteCloser {
				// use delay to create channel baclklog
				// this will demonstrate how Close() is required to drain
				return newWriterWithDelay(closer(output), time.Microsecond)
			},
			func() <-chan time.Time {
				// unused
				return time.Timer{}.C
			},
			bytes.NewBuffer([]byte{}),
		)

		var lastExpectedLog string

		stopCh := make(chan struct{}, 1)

		/*
			CREATE LOTS OF LOGS, RECORDING THE LAST
		*/
		go func() {
			var i int
			for {
				s := time.Now().String()
				lastExpectedLog = s
				l.Print("\n" + s)
				i++

				if i >= numberOfLogs {
					stopCh <- struct{}{}
					return
				}
			}
		}()

		<-stopCh

		if tc.callsClose {
			l.Close()
		}

		b, _ := io.ReadAll(output)

		lines := strings.Split(string(b), "\n")
		lastLine := lines[len(lines)-1]

		// extra line because of new line
		tc.assertNumberOfLines(t, numberOfLogs+1, len(lines))

		tc.assertLastLine(t, lastExpectedLog, lastLine)
	}
}
