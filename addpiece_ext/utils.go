package addpiece_ext

import (
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

var WaitBeforeAddPiece = func() func() {
	var wait_time time.Duration
	if interval := func() int64 {
		interval := os.Getenv("WAIT_BEFORE_ADDPIECE")
		if interval == "" {
			return 0
		}
		num, err := strconv.Atoi(interval)
		if err != nil {
			log.Errorf("[DD] env WAIT_BEFORE_ADDPIECE err:[%v]", err)
			return 0
		}
		return int64(num)
	}(); interval != 0 {
		wait_time = time.Second * time.Duration(interval)
	} else {
		wait_time = time.Second
	}

	var waitForAddPiece int64
	go func() {
		for {
			time.Sleep(time.Second)
			if wait := atomic.LoadInt64(&waitForAddPiece); wait > 0 {
				if wait >= int64(time.Second) {
					atomic.AddInt64(&waitForAddPiece, -int64(time.Second))
				} else {
					atomic.StoreInt64(&waitForAddPiece, 0)
				}
			}
		}
	}()

	return func() {
		n := atomic.AddInt64(&waitForAddPiece, int64(wait_time))
		time.Sleep(time.Duration(n - int64(wait_time)))
		return
	}
}()
