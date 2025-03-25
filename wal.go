package hara

import (
	"context"
	"io"
)

type LogEntry struct {
	Term      uint64
	Offset    uint64
	Value     []byte
	Timestamp uint64
}

// Write a Go interface for a write-ahead log (WAL) that supports concurrent reads and writes.
type Wal interface {
	io.Closer

	Append(ctx context.Context, entry *LogEntry) error

	Sync(ctx context.Context) error

	Truncate(ctx context.Context, offset uint64) error

	LastOffset() uint64

	FirstOffset() uint64
}
