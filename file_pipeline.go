package hara

import (
	"fmt"
	"os"
	"path/filepath"
)

type filePipeline struct {
	dir   string
	size  int64
	count int

	errc  chan error
	donec chan struct{}
	filec chan *os.File
}

func newFilePipeline(dir string, size int64) *filePipeline {
	fp := &filePipeline{
		dir:   dir,
		size:  size,
		count: 0,
		errc:  make(chan error, 1),
		donec: make(chan struct{}),
		filec: make(chan *os.File, 1),
	}

	go fp.run()
	return fp
}

func (p *filePipeline) Open() (f *os.File, err error) {
	select {
	case f = <-p.filec:
	case err = <-p.errc:
	}

	return f, err
}

func (p *filePipeline) Close() error {
	close(p.donec)
	return <-p.errc
}

func (p *filePipeline) run() {
	defer close(p.errc)

	for {
		f, err := p.newFile()
		if err != nil {
			p.errc <- err
			return
		}

		select {
		case p.filec <- f:
		case <-p.donec:
			return
		}
	}
}

func (p *filePipeline) newFile() (*os.File, error) {
	fpath := filepath.Join(p.dir, fmt.Sprintf("%d.tmp", p.count))
	f, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	p.count++
	return f, err
}
