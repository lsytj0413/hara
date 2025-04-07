package hara_test

import (
	"sync"
	"testing"

	"github.com/lsytj0413/hara"
)

func Benchmark_mutex(b *testing.B) {
	v := int64(0)

	lock := sync.Mutex{}
	w := func(i int) {
		lock.Lock()
		v += int64(i)
		lock.Unlock()
	}
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			w(1)
		}
	})
}

func Benchmark_chan(b *testing.B) {
	v := int64(0)

	ch := make(chan int64, 1000)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			i, ok := <-ch
			if ok {
				v += i
				continue
			}

			return
		}
	}()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			ch <- int64(1)
		}
	})
	close(ch)
	wg.Wait()
}

func TestWal(t *testing.T) {
	wal, err := hara.NewWal("/tmp")
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 1000; i++ {
		err = wal.Append(nil, &hara.Entry{
			Offset: uint64(i),
			Data:   []byte("hello world"),
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	err = wal.Sync(nil)
	if err != nil {
		t.Fatal(err)
	}

	err = wal.Truncate(nil, 0)
	if err != nil {
		t.Fatal(err)
	}

	err = wal.Close()
	if err != nil {
		t.Fatal(err)
	}
}
