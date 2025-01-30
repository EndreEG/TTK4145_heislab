package elevator

import (
	"fmt"
)

//Imported
const NumFloors int = 4
const NumButtons int = 3

//Elevator_io.h
type Dirn int

const (
	D_Down Dirn = -1
	D_Stop Dirn = 0
	D_Up   Dirn = 1
)

type Button int

const (
	B_HallUp Button = iota
	B_Halldown
	B_Cab
)

//----------------------------

//Elevator.h
type ElevatorBehaviour int

const (
	EB_Idle = iota
	EB_DoorOpen
	EB_Moving
)

type ClearRequestVariant int

const (
	CV_All = iota
	CV_InDirn
)

type Config struct {
	clearRequestVariant ClearRequestVariant
	doorOpenDuration_s  float64
}

type Elevator struct {
	floor     int
	dirn      Dirn
	request   [NumFloors][NumButtons]int
	behaviour ElevatorBehaviour
	config    Config
}

//------------------------------------

//elevator.c
func eb_toString(eb ElevatorBehaviour) string {
	if eb == EB_Idle {
		return "EB_Idle"
	} else if eb == EB_DoorOpen {
		return "EB_DoorOpen"
	} else if eb == EB_Moving {
		return "EB_Moving"
	} else {
		return "EB_UNDEFINED"
	}
}

func elevio_dirn_toString(dr Dirn) string { //Fant aldri i c-fil, men logikken burde være den samme som forrige
	if dr == D_Down {
		return "D_Down"
	} else if dr == D_Stop {
		return "D_Stop"
	} else if dr == D_Up {
		return "D_Up"
	} else {
		return "DR_UNDEFINED"
	}
}



func elevator_uninitialized() Elevator {
	el := Elevator{floor: -1,
		dirn:      D_Stop,
		behaviour: EB_Idle,
		config:    Config{clearRequestVariant: CV_All, doorOpenDuration_s: 3.0},
	}
	return el
}

func elevator_print(es Elevator) {
	fmt.Printf("  +--------------------+\n")
	fmt.Printf("|  floor = %-15d|\n|  dirn = %-15s |\n|  behav = %-15s|\n", es.floor,
		elevio_dirn_toString(es.dirn), eb_toString(es.behaviour))
	fmt.Printf("  +--------------------+\n")
	fmt.Printf("|   | up  | dn  | cab |   |\n")
	for f := NumFloors - 1; f >= 0; f-- {
		fmt.Printf("|  %d", f)
		for btn := 0; btn < NumButtons; btn++ {
			if ((f == NumFloors-1) && (btn == int(B_HallUp))) || ((f == 0) && (btn == int(B_Halldown))) {
				fmt.Printf("|     ")
			} else {
				if es.request[f][btn] == 1 {
					fmt.Printf("|  #  ")
				} else {
					fmt.Printf("|  -  ")
				}

			}

		}

		fmt.Printf("|   |\n")
	}
	fmt.Printf("  +--------------------+\n")
}

func main() {
	a := elevator_uninitialized()
	elevator_print(a)
}
