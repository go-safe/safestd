package safestd

import (
	"io"
	"log"
	"os"
	"runtime"
	"sync/atomic"
)

var isFinalizerEnabled int32 = 1

func IsCheckEnabled() bool {
	return atomic.LoadInt32(&isFinalizerEnabled) == 1
}

func SetRuntimeChecks(enabled bool) {
	var value int32 = 0
	if enabled {
		value = 1
	}
	atomic.StoreInt32(&isFinalizerEnabled, value)
}

type File struct {
	finalizerRate int64
	log           io.Reader
}

func (f *File) Open(name string) (*os.File, CancelFn, error) {
	return nil, nil, nil
}

func OpenFile(name string) (*os.File, CancelFn, error) {
	file, openErr := os.Open(name)
	if openErr != nil {
		return nil, nil, openErr
	}

	c := &canceler{file: file}
	if IsCheckEnabled() {
		runtime.SetFinalizer(c, checkCanceler)
	}

	return file, c.Close, nil
}

type canceler struct {
	file  *os.File
	state int32
}

func (c *canceler) Close() error {
	atomic.StoreInt32(&c.state, 1)
	return c.file.Close()
}

func checkCanceler(c *canceler) {
	if atomic.LoadInt32(&c.state) == 1 {
		return
	}
	log.Printf("file isn't closed properly: %#v\n", c.file.Name())
}
