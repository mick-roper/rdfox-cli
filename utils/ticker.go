package utils

import "time"

func DoWithTicker(action func(), onTick func()) {
	tick := time.Tick(time.Second * 1)
	stop := make(chan struct{})
	defer close(stop)

	go func() {
		for {
			select {
			case <-tick:
				onTick()
			case <-stop:
				return
			}
		}
	}()

	action()
}
