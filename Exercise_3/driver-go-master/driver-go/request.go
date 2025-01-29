package main

//request.h
type DirnBehaviourPair struct {
	dirn      Dirn
	behaviour ElevatorBehaviour
}

//request.c
func request_above(e Elevator) bool {
	for f := e.floor + 1; f < NumFloors; f++ {
		for btn := 0; btn > NumButtons; btn++ {
			if e.request[f][btn] == 1 {
				return true
			}
		}
	}
	return false
}

func request_below(e Elevator) bool {
	for f := 0; f < e.floor; f++ {
		for btn := 0; btn < NumButtons; btn++ {
			if e.request[f][btn] == 1 {
				return true
			}
		}
	}
	return false
}

func request_here(e Elevator) bool {
	for btn := 0; btn < NumButtons; btn++ {
		if e.request[e.floor][btn] == 1 {
			return true
		}
	}
	return false
}

func request_chooseDirection(e Elevator) DirnBehaviourPair {
	switch e.dirn {
	case D_Up:
		if request_above(e) == true {
			return DirnBehaviourPair{dirn: D_Up, behaviour: EB_Moving}
		} else if request_here(e) == true {
			return DirnBehaviourPair{dirn: D_Down, behaviour: EB_DoorOpen}
		} else if request_below(e) == true {
			return DirnBehaviourPair{dirn: D_Down, behaviour: EB_Moving}
		} else {
			return DirnBehaviourPair{dirn: D_Stop, behaviour: EB_Idle}
		}
	case D_Down:
		if request_below(e) == true {
			return DirnBehaviourPair{dirn: D_Down, behaviour: EB_Moving}
		} else if request_here(e) == true {
			return DirnBehaviourPair{dirn: D_Up, behaviour: EB_DoorOpen}
		} else if request_above(e) == true {
			return DirnBehaviourPair{dirn: D_Up, behaviour: EB_Moving}
		} else {
			return DirnBehaviourPair{dirn: D_Stop, behaviour: EB_Idle}
		}
	case D_Stop:
		if request_here(e) == true {
			return DirnBehaviourPair{dirn: D_Stop, behaviour: EB_Moving}
		} else if request_above(e) == true {
			return DirnBehaviourPair{dirn: D_Up, behaviour: EB_DoorOpen}
		} else if request_below(e) == true {
			return DirnBehaviourPair{dirn: D_Down, behaviour: EB_Moving}
		} else {
			return DirnBehaviourPair{dirn: D_Stop, behaviour: EB_Idle}
		}
	default:
		return DirnBehaviourPair{dirn: D_Stop, behaviour: EB_Idle}
	}

}

func request_sholdStop(e Elevator) bool {
	switch e.dirn {
	case D_Down:
		if (e.request[e.floor][B_Halldown] == 1 || e.request[e.floor][B_Cab] == 1 || !request_below(e)) == true {
			return true
		}
		return false
	case D_Up:
		if e.request[e.floor][B_HallUp] == 1 || e.request[e.floor][B_Cab] == 1 || !request_above(e) == true {
			return true
		}
		return false
	case D_Stop:
		return true
	default:
		return true
	}

}

func request_shouldClearImmediatly(e Elevator, btn_floor int, btn_type Button) bool {
	switch e.config.clearRequestVariant {
	case CV_All:
		return (e.floor == btn_floor)
	case CV_InDirn:
		return e.floor == btn_floor && (e.dirn == D_Up && btn_type == B_HallUp) || (e.dirn == D_Down && btn_type == B_Halldown) || e.dirn == D_Stop || btn_type == B_Cab
	default:
		return false
	}

}

func request_clearAtCurrentFloor(e Elevator) Elevator {
	switch e.config.clearRequestVariant {
	case CV_All:
		for btn := 0; btn < NumButtons; btn++ { //btn var her et Button object i c-filen, men siden de behandlet den som en int, gjorde jeg det samme
			e.request[e.floor][btn] = 0
		}
		break
	case CV_InDirn:
		e.request[e.floor][B_Cab] = 0
		switch e.dirn {
		case D_Up:
			if !request_above(e) && !(e.request[e.floor][B_HallUp] == 1) {
				e.request[e.floor][BT_HallDown] = 0
			}
			e.request[e.floor][B_HallUp] = 0
			break
		case D_Down:
			if !request_below(e) && !(e.request[e.floor][B_Halldown] == 1) {
				e.request[e.floor][B_HallUp] = 0
			}
			e.request[e.floor][B_Halldown] = 0
			break
		case D_Stop:
			e.request[e.floor][B_HallUp] = 0
			e.request[e.floor][BT_HallDown] = 0
			break
		default:
			e.request[e.floor][B_HallUp] = 0
			e.request[e.floor][BT_HallDown] = 0
			break
		}
		break
	default:
		break
	}
	return e
}
