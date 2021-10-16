package main

import "time"

type Ticker struct {
	ticks    time.Duration
	interval time.Duration
}

func (ticker *Ticker) Initialize(interval time.Duration) {
	ticker.ticks = 0
	ticker.interval = interval
}

func (ticker *Ticker) Advance(timePassed time.Duration) bool {
	ticker.ticks += timePassed
	triggered := ticker.ticks > ticker.interval
	if triggered {
		ticker.ticks = 0
	}
	return triggered
}
