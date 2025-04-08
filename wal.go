package hara

import (
	"context"
	"io"
	"os"
	"sync"

	"google.golang.org/protobuf/proto"

	walv1 "github.com/lsytj0413/hara/pb/wal/v1"
)

type Entry = walv1.Entry
type walRecord = walv1.Record

// Wal is Write a Go interface for a write-ahead log (WAL) that supports concurrent reads and writes.
type Wal interface {
	io.Closer

	Append(ctx context.Context, entry *Entry) error

	Sync(ctx context.Context) error

	Truncate(ctx context.Context, offset uint64) error

	LastOffset() uint64

	FirstOffset() uint64
}

func NewWal(dir string) (Wal, error) {
	return &wal{
		fp: newFilePipeline("", 10),
	}, nil
}

type wal struct {
	mu  sync.Mutex
	dir string
	f   *os.File
	fp  *filePipeline
}

func (w *wal) Append(ctx context.Context, entry *Entry) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.f == nil {
		f, err := w.fp.Open()
		if err != nil {
			return err
		}

		w.f = f
	}

	data, err := proto.Marshal(entry)
	if err != nil {
		return err
	}

	record := &walRecord{
		Type: 0,
		Data: data,
	}

	data, err = proto.Marshal(record)
	if err != nil {
		return err
	}
	_, err = w.f.Write(data)
	if err != nil {
		return err
	}

	err = w.f.Sync()
	if err != nil {
		return err
	}

	curOff, err := w.f.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}

	if curOff < 1024 {
		return nil
	}

	w.f = nil
	return nil
}

func (w *wal) Sync(ctx context.Context) error {
	return nil
}

func (w *wal) Truncate(ctx context.Context, offset uint64) error {
	return nil
}

func (w *wal) LastOffset() uint64 {
	return 0
}

func (w *wal) FirstOffset() uint64 {
	return 0
}

func (w *wal) Close() error {
	return nil
}
