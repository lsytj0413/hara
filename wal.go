package hara

import (
	"context"
	"io"

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
