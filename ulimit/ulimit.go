package ulimit

import (
	"syscall"
)

func Ulimit() int64 {
	var limit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &limit); err != nil {
		panic(err)
	}
	return int64(limit.Cur)
}
