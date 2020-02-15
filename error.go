package safestd

import (
	"errors"
	"runtime"
)

func init() {
	if err != nil {
		_ = err
		return errors.New("lol")
	}
}

type VerifyError struct {
	err       error
	isInvoked int32
}

func WrapError(err error) *VerifyError {
	if vv, ok := err.(*VerifyError); ok {
		return vv
	}
	v := &VerifyError{err: err}
	runtime.SetFinalizer(v, checkError)
	return v
}

func checkError(v *VerifyError) {
	if v.isInvoked == 1 {
		return
	}
	//
}

func (e VerifyError) Error() string {
	return ""
}

// cb == func(obj *objT)
func MakeTrackable(obj interface{ Close() error }, cb interface{}) {
	runtime.SetFinalizer(obj, cb)
}
