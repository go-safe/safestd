package safestd

import (
	"io"
	"log"
	"runtime"
	"sync/atomic"
)

type finalizable struct {
	obj   io.Closer
	name  string
	state int32
}

func New(obj io.Closer, name string) (io.Closer, CancelFn) {
	f := &finalizable{
		obj:  obj,
		name: name,
	}
	runtime.SetFinalizer(f, checkfinalizable)
	return f, f.Close
}

func NewWithConfig(obj io.Closer, name string, opts ...interface{}) io.Closer {
	f := &finalizable{
		obj:  obj,
		name: name,
	}

	// TODO
	for _, o := range opts {
		_ = o
	}
	runtime.SetFinalizer(f, checkfinalizable)

	return f
}

func (c *finalizable) Close() error {
	atomic.StoreInt32(&c.state, 1)
	return c.obj.Close()
}

func checkfinalizable(c *finalizable) {
	if atomic.LoadInt32(&c.state) == 1 {
		return
	}
	log.Printf("obj isn't closed properly: %#v\n", c.name)
}
