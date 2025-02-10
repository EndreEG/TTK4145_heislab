package main

import (
	"Heis/elevio"
	"fmt"
)

// Imported
const NumFloors int = 4
const NumButtons int = 3

//----------------------------

// Elevator.h
type ElevatorBehaviour int

const (
	EB_Idle = iota
	EB_DoorOpen
	EB_Moving
	EB_Stop
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
	dirn      elevio.MotorDirection
	request   [NumFloors][NumButtons]int
	behaviour ElevatorBehaviour
	config    Config
}

//------------------------------------

// elevator.c
func Elevio_toString(eb ElevatorBehaviour) string {
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

func Elevio_button_toString(eb elevio.ButtonType) string {
	if eb == elevio.BT_HallUp {
		return "BT_HallUp"
	} else if eb == elevio.BT_Cab {
		return "BT_HallDown"
	} else if eb == elevio.BT_HallDown {
		return "BT_HallDown"
	} else {
		return "BT_UNDEFINED"
	}
}

func Elevio_dirn_toString(dr elevio.MotorDirection) string { //Fant aldri i c-fil, men logikken burde være den samme som forrige
	if dr == elevio.MD_Down {
		return "D_Down"
	} else if dr == elevio.MD_Stop {
		return "D_Stop"
	} else if dr == elevio.MD_Up {
		return "D_Up"
	} else {
		return "DR_UNDEFINED"
	}
}

func Elevator_print(es Elevator) {
	fmt.Printf("  +--------------------+\n")
	fmt.Printf("|  floor = %-15d|\n|  dirn = %-15s |\n|  behav = %-15s|\n", es.floor,
		Elevio_dirn_toString(es.dirn), Elevio_toString(es.behaviour))
	fmt.Printf("  +--------------------+\n")
	fmt.Printf("|   | up  | dn  | cab |   |\n")
	for f := NumFloors - 1; f >= 0; f-- {
		fmt.Printf("|  %d", f)
		for btn := 0; btn < NumButtons; btn++ {
			if ((f == NumFloors-1) && (btn == int(elevio.BT_HallUp))) || ((f == 0) && (btn == int(elevio.BT_HallDown))) {
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

func Elevator_uninitialized() Elevator {
	el := Elevator{floor: -1,
		dirn:      elevio.MD_Stop,
		behaviour: EB_Idle,
		config:    Config{clearRequestVariant: CV_All, doorOpenDuration_s: 3.0},
	}
	return el
}
