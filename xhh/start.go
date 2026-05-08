package xhh

import "time"

func Start() {
	go func() {
		for {
			AutoReply()
			time.Sleep(5 * time.Second)
			CheckAt()
			time.Sleep(5 * time.Second)
		}
	}()
}
