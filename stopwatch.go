package main

import (
	"time"
)

type StopWatch struct {
	start, stop time.Time
}

func Start() time.Time {
	return time.Now()
}

func Stop(start time.Time) *StopWatch {
	watch := StopWatch{start: start, stop: time.Now()}
	return &watch
}

func (self *StopWatch) Milliseconds() uint32 {
	return uint32(self.stop.Sub(self.start) / time.Millisecond)
}
