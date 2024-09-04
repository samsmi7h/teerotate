package teerotate

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"
)

type makeNewOutput func() WriteCloseSizer

const dateFmt = "2006-01-02T15:04:05"

type Sizer interface {
	SizeInBytes() ByteSize
}

type WriteCloseSizer interface {
	io.WriteCloser
	Sizer
}

func fileFactory(dir string) makeNewOutput {
	return func() WriteCloseSizer {
		t := time.Now().Format(dateFmt)
		closedFileName := fmt.Sprintf("%s.log", t)
		closedPath := path.Join(dir, closedFileName)
		livePath := fmt.Sprintf("%s.live", closedPath)

		f, err := os.Create(livePath)
		must(err)

		return &renameOnCloseFile{
			File:           f,
			closedFilePath: closedPath,
		}
	}
}

type renameOnCloseFile struct {
	*os.File

	closedFilePath string
}

func (r *renameOnCloseFile) Close() error {
	if err := r.File.Close(); err != nil {
		return err
	}

	fmt.Println("renaming closed file...", r.File.Name(), r.closedFilePath)
	if err := os.Rename(r.File.Name(), r.closedFilePath); err != nil {
		panic(err)
	}

	return nil
}

func (r *renameOnCloseFile) SizeInBytes() ByteSize {
	s, err := r.File.Stat()
	if err != nil {
		return 0
	}

	if s == nil {
		return 0
	}

	return ByteSize(s.Size())
}
