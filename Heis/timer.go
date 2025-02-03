package main

import (
	"time"
)

const _pollRate = 20 * time.Millisecond

func GetWallTime() float64 { //Får verdenstiden i sekunder.
	tid := time.Now()
	return float64(tid.UnixNano()) / 1e9
}

var timerEndTime float64
var timerActive int

func Timer_start(duration float64) {
	timerEndTime = GetWallTime() + duration
	timerActive = 1
}

func Timer_stop() {
	timerActive = 0
}

func Timer_timedOut() bool {
	if timerActive == 1 && (GetWallTime() > timerEndTime) {
		return true
	} else {
		return false
	}
}

func PollTimer(reciver chan<- bool){
	prev := false
	for {
		time.Sleep(_pollRate)
		v := Timer_timedOut()
		if v != prev {
			reciver <- v
		}
		prev = v
	}
}
