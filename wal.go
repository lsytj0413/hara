package hara

import (
	"context"
	"io"
	"os"
	"path"

	"google.golang.org/protobuf/proto"

	"github.com/lsytj0413/hara/pb"
)

// Wal is Write a Go interface for a write-ahead log (WAL) that supports concurrent reads and writes.
type Wal interface {
	io.Closer

	Append(ctx context.Context, entry *pb.Entry) error

	Sync(ctx context.Context) error

	Truncate(ctx context.Context, offset uint64) error

	LastOffset() uint64

	FirstOffset() uint64
}

func NewWal(dir string) (Wal, error) {
	return &wal{}, nil
}

type wal struct {
	dir string
	f   *os.File
}

func (w *wal) Append(ctx context.Context, entry *pb.Entry) error {
	if w.f == nil {
		f, err := os.OpenFile(path.Join(w.dir, "0.wal"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return err
		}

		w.f = f
	}

	data, err := proto.Marshal(entry)
	if err != nil {
		return err
	}

	record := &pb.Record{
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
