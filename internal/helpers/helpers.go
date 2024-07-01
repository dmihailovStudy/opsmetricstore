package helpers

import "time"

func Wait(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}
