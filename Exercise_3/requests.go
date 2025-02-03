package elevator

const (
	N_FLOORS  = 4
	N_BUTTONS = 3
)

type Direction int
type Button int
type ClearRequestVariant int
type ElevatorBehaviour int

const (
	D_Up Direction = iota
	D_Down
	D_Stop
)

const (
	B_HallUp Button = iota
	B_HallDown
	B_Cab
)

const (
	CV_All ClearRequestVariant = iota
	CV_InDirn
)

const (
	EB_Moving ElevatorBehaviour = iota
	EB_DoorOpen
	EB_Idle
)

type DirnBehaviourPair struct {
	Dirn      Direction
	Behaviour ElevatorBehaviour
}

type Config struct {
	ClearRequestVariant ClearRequestVariant
}

type Elevator struct {
	Floor    int
	Dirn     Direction
	Requests [N_FLOORS][N_BUTTONS]bool
	Config   Config
}

func requestsAbove(e Elevator) bool {
	for f := e.Floor + 1; f < N_FLOORS; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if e.Requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func requestsBelow(e Elevator) bool {
	for f := 0; f < e.Floor; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if e.Requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func requestsHere(e Elevator) bool {
	for btn := 0; btn < N_BUTTONS; btn++ {
		if e.Requests[e.Floor][btn] {
			return true
		}
	}
	return false
}

func requestsChooseDirection(e Elevator) DirnBehaviourPair {
	switch e.Dirn {
	case D_Up:
		if requestsAbove(e) {
			return DirnBehaviourPair{D_Up, EB_Moving}
		} else if requestsHere(e) {
			return DirnBehaviourPair{D_Down, EB_DoorOpen}
		} else if requestsBelow(e) {
			return DirnBehaviourPair{D_Down, EB_Moving}
		}
	case D_Down:
		if requestsBelow(e) {
			return DirnBehaviourPair{D_Down, EB_Moving}
		} else if requestsHere(e) {
			return DirnBehaviourPair{D_Up, EB_DoorOpen}
		} else if requestsAbove(e) {
			return DirnBehaviourPair{D_Up, EB_Moving}
		}
	case D_Stop:
		if requestsHere(e) {
			return DirnBehaviourPair{D_Stop, EB_DoorOpen}
		} else if requestsAbove(e) {
			return DirnBehaviourPair{D_Up, EB_Moving}
		} else if requestsBelow(e) {
			return DirnBehaviourPair{D_Down, EB_Moving}
		}
	}
	return DirnBehaviourPair{D_Stop, EB_Idle}
}

func requestsShouldStop(e Elevator) bool {
	switch e.Dirn {
	case D_Down:
		return e.Requests[e.Floor][B_HallDown] || e.Requests[e.Floor][B_Cab] || !requestsBelow(e)
	case D_Up:
		return e.Requests[e.Floor][B_HallUp] || e.Requests[e.Floor][B_Cab] || !requestsAbove(e)
	case D_Stop:
		return true
	}
	return false
}

func requestsShouldClearImmediately(e Elevator, btnFloor int, btnType Button) bool {
	switch e.Config.ClearRequestVariant {
	case CV_All:
		return e.Floor == btnFloor
	case CV_InDirn:
		return e.Floor == btnFloor && (e.Dirn == D_Up && btnType == B_HallUp ||
			e.Dirn == D_Down && btnType == B_HallDown || e.Dirn == D_Stop || btnType == B_Cab)
	}
	return false
}

func requestsClearAtCurrentFloor(e Elevator) Elevator {
	switch e.Config.ClearRequestVariant {
	case CV_All:
		for btn := 0; btn < N_BUTTONS; btn++ {
			e.Requests[e.Floor][btn] = false
		}
	case CV_InDirn:
		e.Requests[e.Floor][B_Cab] = false
		switch e.Dirn {
		case D_Up:
			if !requestsAbove(e) && !e.Requests[e.Floor][B_HallUp] {
				e.Requests[e.Floor][B_HallDown] = false
			}
			e.Requests[e.Floor][B_HallUp] = false
		case D_Down:
			if !requestsBelow(e) && !e.Requests[e.Floor][B_HallDown] {
				e.Requests[e.Floor][B_HallUp] = false
			}
			e.Requests[e.Floor][B_HallDown] = false
		case D_Stop:
			e.Requests[e.Floor][B_HallUp] = false
			e.Requests[e.Floor][B_HallDown] = false
		}
	}
	return e
}
