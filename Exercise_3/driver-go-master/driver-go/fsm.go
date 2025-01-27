package fsm

import (
	"Driver-go/elevio"
	"fmt"
	"time"
)

// Elevator constants
const numFloors = 4

// Elevator behavior states
type ElevatorBehavior int

const (
	EB_Idle ElevatorBehavior = iota
	EB_Moving
	EB_DoorOpen
)

// Direction constants
type MotorDirection elevio.MotorDirection

const (
	D_Up   = MotorDirection(elevio.MD_Up)
	D_Down = MotorDirection(elevio.MD_Down)
	D_Stop = MotorDirection(elevio.MD_Stop)
)

// Elevator struct
type Elevator struct {
	Floor        int
	Dirn         MotorDirection
	Behaviour    ElevatorBehavior
	Requests     [numFloors][3]bool // Requests[floor][buttonType]
	DoorOpenTime time.Duration      // Door open duration (in seconds)
}

// Global elevator instance
var elevator Elevator

// Initialize the FSM
func Init() {
	elevator = Elevator{
		Floor:        -1, // Unknown floor
		Dirn:         D_Stop,
		Behaviour:    EB_Idle,
		Requests:     [numFloors][3]bool{},
		DoorOpenTime: 3 * time.Second, // Door stays open for 3 seconds
	}
	fmt.Println("FSM initialized.")
}

// Update all button lights based on requests
func updateAllLights() {
	for floor := 0; floor < numFloors; floor++ {
		for btn := elevio.BT_HallUp; btn <= elevio.BT_Cab; btn++ {
			elevio.SetButtonLamp(btn, floor, elevator.Requests[floor][btn])
		}
	}
}

// Handle request button press
func OnRequestButtonPress(btnFloor int, btnType elevio.ButtonType) {
	fmt.Printf("\nButton pressed: floor=%d, type=%v\n", btnFloor, btnType)

	elevator.Requests[btnFloor][btnType] = true

	switch elevator.Behaviour {
	case EB_Idle:
		startMovingIfIdle()

	case EB_Moving:
		// Requests are simply queued until the elevator reaches the target floor.

	case EB_DoorOpen:
		if elevator.Floor == btnFloor {
			fmt.Println("Restarting door timer.")
		} else {
			elevator.Requests[btnFloor][btnType] = true
		}
	}

	updateAllLights()
}

// Handle floor arrival
func OnFloorArrival(newFloor int) {
	fmt.Printf("\nArrived at floor: %d\n", newFloor)

	elevator.Floor = newFloor
	elevio.SetFloorIndicator(newFloor)

	if elevator.Behaviour == EB_Moving && shouldStop() {
		elevio.SetMotorDirection(elevio.MD_Stop)
		elevator.Behaviour = EB_DoorOpen
		clearRequestsAtCurrentFloor()
		elevio.SetDoorOpenLamp(true)
		time.AfterFunc(elevator.DoorOpenTime, func() {
			OnDoorTimeout()
		})
	}

	updateAllLights()
}

// Handle door timeout
func OnDoorTimeout() {
	fmt.Println("\nDoor timeout.")

	if elevator.Behaviour == EB_DoorOpen {
		elevio.SetDoorOpenLamp(false)
		startMovingIfIdle()
	}
}

// Determine if the elevator should stop at the current floor
func shouldStop() bool {
	if elevator.Dirn == D_Stop {
		return false
	}
	return elevator.Requests[elevator.Floor][elevio.BT_Cab] ||
		(elevator.Dirn == D_Up && elevator.Requests[elevator.Floor][elevio.BT_HallUp]) ||
		(elevator.Dirn == D_Down && elevator.Requests[elevator.Floor][elevio.BT_HallDown])
}

// Clear requests at the current floor
func clearRequestsAtCurrentFloor() {
	for btn := elevio.BT_HallUp; btn <= elevio.BT_Cab; btn++ {
		elevator.Requests[elevator.Floor][btn] = false
		elevio.SetButtonLamp(btn, elevator.Floor, false)
	}
}

// Start moving if idle
func startMovingIfIdle() {
	nextDir := chooseDirection()
	if nextDir == D_Stop {
		elevator.Behaviour = EB_Idle
	} else {
		elevator.Dirn = nextDir
		elevio.SetMotorDirection(elevio.MotorDirection(nextDir))
		elevator.Behaviour = EB_Moving
	}
}

// Choose the direction based on pending requests
func chooseDirection() MotorDirection {
	if elevator.Dirn == D_Up || elevator.Dirn == D_Stop {
		for f := elevator.Floor + 1; f < numFloors; f++ {
			for btn := elevio.BT_HallUp; btn <= elevio.BT_Cab; btn++ {
				if elevator.Requests[f][btn] {
					return D_Up
				}
			}
		}
	}
	if elevator.Dirn == D_Down || elevator.Dirn == D_Stop {
		for f := elevator.Floor - 1; f >= 0; f-- {
			for btn := elevio.BT_HallUp; btn <= elevio.BT_Cab; btn++ {
				if elevator.Requests[f][btn] {
					return D_Down
				}
			}
		}
	}
	return D_Stop
}
