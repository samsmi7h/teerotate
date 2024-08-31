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

const dateFmt = "2006-01-02T15:04:05"

func fileFactory(dir string) makeNewOutput {
	return func() io.WriteCloser {
		fileName := fmt.Sprintf("%s_%s.log", uuid.NewString(), time.Now().Format(dateFmt))
		p := path.Join(dir, fileName)
		f, err := os.Create(p)
		must(err)

		return f
	}
}
