package fsm


type ElevatorState int

const (
	Idle ElevatorState = iota
	Initializing
	GoingToFloor
	AtFloor
	Emergency
)



func updateState()