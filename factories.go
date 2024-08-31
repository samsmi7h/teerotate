package teerotate

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/google/uuid"
)

type makeNewTicker func() <-chan time.Time

func tickerFactory(lifespan time.Duration) makeNewTicker {
	return func() <-chan time.Time {
		return time.NewTicker(lifespan).C
	}
}

type makeNewOutput func() io.WriteCloser

func fileFactory(dir string) makeNewOutput {
	return func() io.WriteCloser {
		fileName := fmt.Sprintf("%s_%s.log", uuid.NewString(), time.Now().Format(time.DateTime))
		p := path.Join(dir, fileName)
		f, err := os.Create(p)
		must(err)

		return f
	}
}
