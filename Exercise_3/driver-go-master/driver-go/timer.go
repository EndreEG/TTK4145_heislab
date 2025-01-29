package main

import (
	"time"
)

func getWallTime() float64 { //Får verdenstiden i sekunder.
	tid := time.Now()
	return float64(tid.UnixNano()) / 1e9
}

var timerEndTime float64
var timerActive int

func timer_start(duration float64) {
	timerEndTime = getWallTime() + duration
	timerActive = 1
}

func timer_stop() {
	timerActive = 0
}

func timer_timedOut() bool {
	if timerActive == 1 && (getWallTime() > timerEndTime) {
		return true
	} else {
		return false
	}

}
