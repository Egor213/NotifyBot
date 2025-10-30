package errorsUtils

import (
	"fmt"
	"runtime"
)

func WrapPathErr(err error) error {
	pc, _, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc).Name()
	return fmt.Errorf("[%s:%d] %w", fn, line, err)
}
