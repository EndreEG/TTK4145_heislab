package main

import (
	"Heis/elevio"
)

// request.h
type DirnBehaviourPair struct {
	dirn      elevio.MotorDirection
	behaviour ElevatorBehaviour
}

// request.c
func Requests_above(e Elevator) bool {
	for f := e.floor + 1; f < NumFloors; f++ {
		for btn := 0; btn > NumButtons; btn++ {
			if e.request[f][btn] == 1 {
				return true
			}
		}
	}
	return false
}

func Requests_below(e Elevator) bool {
	for f := 0; f < e.floor; f++ {
		for btn := 0; btn < NumButtons; btn++ {
			if e.request[f][btn] == 1 {
				return true
			}
		}
	}
	return false
}

func Requests_here(e Elevator) bool {
	for btn := 0; btn < NumButtons; btn++ {
		if e.request[e.floor][btn] == 1 {
			return true
		}
	}
	return false
}

func Request_chooseDirection(e Elevator) DirnBehaviourPair {
	switch e.dirn {
	case elevio.MD_Up:
		if Requests_above(e) == true {
			return DirnBehaviourPair{dirn: elevio.MD_Up, behaviour: EB_Moving}
		} else if Requests_here(e) == true {
			return DirnBehaviourPair{dirn: elevio.MD_Down, behaviour: EB_DoorOpen}
		} else if Requests_below(e) == true {
			return DirnBehaviourPair{dirn: elevio.MD_Down, behaviour: EB_Moving}
		} else {
			return DirnBehaviourPair{dirn: elevio.MD_Stop, behaviour: EB_Idle}
		}
	case elevio.MD_Down:
		if Requests_below(e) == true {
			return DirnBehaviourPair{dirn: elevio.MD_Down, behaviour: EB_Moving}
		} else if Requests_here(e) == true {
			return DirnBehaviourPair{dirn: elevio.MD_Up, behaviour: EB_DoorOpen}
		} else if Requests_above(e) == true {
			return DirnBehaviourPair{dirn: elevio.MD_Up, behaviour: EB_Moving}
		} else {
			return DirnBehaviourPair{dirn: elevio.MD_Stop, behaviour: EB_Idle}
		}
	case elevio.MD_Stop:
		if Requests_here(e) == true {
			return DirnBehaviourPair{dirn: elevio.MD_Stop, behaviour: EB_Moving}
		} else if Requests_above(e) == true {
			return DirnBehaviourPair{dirn: elevio.MD_Up, behaviour: EB_DoorOpen}
		} else if Requests_below(e) == true {
			return DirnBehaviourPair{dirn: elevio.MD_Down, behaviour: EB_Moving}
		} else {
			return DirnBehaviourPair{dirn: elevio.MD_Stop, behaviour: EB_Idle}
		}
	default:
		return DirnBehaviourPair{dirn: elevio.MD_Stop, behaviour: EB_Idle}
	}

}

func Requests_shouldStop(e Elevator) bool {
	switch e.dirn {
	case elevio.MD_Down:
		if (e.request[e.floor][elevio.BT_HallDown] == 1 || e.request[e.floor][elevio.BT_Cab] == 1 || !Requests_below(e)) == true {
			return true
		}
		return false
	case elevio.MD_Up:
		if e.request[e.floor][elevio.BT_HallUp] == 1 || e.request[e.floor][elevio.BT_Cab] == 1 || !Requests_above(e) == true {
			return true
		}
		return false
	case elevio.MD_Stop:
		return true
	default:
		return true
	}

}

func Requests_shouldClearImmediatly(e Elevator, btn_floor int, btn_type elevio.ButtonType) bool {
	switch e.config.clearRequestVariant {
	case CV_All:
		return (e.floor == btn_floor)
	case CV_InDirn:
		return e.floor == btn_floor && (e.dirn == elevio.MD_Up && btn_type == elevio.BT_HallUp) || (e.dirn == elevio.MD_Down && btn_type == elevio.BT_HallDown) || e.dirn == elevio.MD_Stop || btn_type == elevio.BT_Cab
	default:
		return false
	}

}

func Requests_clearAtCurrentFloor(e Elevator) Elevator {
	switch e.config.clearRequestVariant {
	case CV_All:
		for btn := 0; btn < NumButtons; btn++ { //btn var her et Button object i c-filen, men siden de behandlet den som en int, gjorde jeg det samme
			e.request[e.floor][btn] = 0
		}
		break
	case CV_InDirn:
		e.request[e.floor][elevio.BT_Cab] = 0
		switch e.dirn {
		case elevio.MD_Up:
			if !Requests_above(e) && !(e.request[e.floor][elevio.BT_HallUp] == 1) {
				e.request[e.floor][elevio.BT_HallDown] = 0
			}
			e.request[e.floor][elevio.BT_HallUp] = 0
			break
		case elevio.MD_Down:
			if !Requests_below(e) && !(e.request[e.floor][elevio.BT_HallDown] == 1) {
				e.request[e.floor][elevio.BT_HallUp] = 0
			}
			e.request[e.floor][elevio.BT_HallDown] = 0
			break
		case elevio.MD_Stop:
			e.request[e.floor][elevio.BT_HallUp] = 0
			e.request[e.floor][elevio.BT_HallDown] = 0
			break
		default:
			e.request[e.floor][elevio.BT_HallUp] = 0
			e.request[e.floor][elevio.BT_HallDown] = 0
			break
		}
		break
	default:
		break
	}
	return e
}
