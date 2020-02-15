package safestd_test

import (
	"log"
	"os"
	"runtime"
	"testing"

	"github.com/safe-go/safestd"
)

func TestOpenFile(t *testing.T) {
	doit()

	runtime.GC()
}

func foo() {
	f, err := os.Open("bombom.txt")

	// 1.
	// f, err := safestd.Open("bombom.txt")
	//
	// 2.
	// fm := safestd.NewFileManager()
	//
	// f, cancel, err := fm.Open("bombom.txt")

	_ = f
	_ = err
}

type App struct {
	fm fileManager
}

type fileManager interface {
	Open(name string) (*os.File, func() error, error)
}

func doit() {
	file, cancel, openErr := safestd.OpenFile("os.go")
	if openErr != nil {
		return
	}
	_ = file
	_ = cancel
}

func wrapClose(cancel func() error) func() {
	return func() {
		if err := cancel(); err != nil {
			log.Printf("omg: %#v", err)
		}
	}
}
