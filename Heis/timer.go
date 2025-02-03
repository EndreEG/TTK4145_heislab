package main

import (
	"time"
)

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
