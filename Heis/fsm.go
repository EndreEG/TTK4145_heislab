package main

import (
	"Heis/elevio"
	"fmt"
	"runtime"
)

var elevator Elevator = Elevator_uninitialized()

func SetAllLights(es Elevator) {
	bu := []elevio.ButtonType{elevio.BT_HallUp, elevio.BT_HallDown, elevio.BT_Cab}
	for floor := 0; floor < NumFloors; floor++ {
		for _, B := range bu {
			elevio.SetButtonLamp(B, floor, es.request[floor][B] == 1)
		}
	}
}

func Fsm_onInitBetweenFloors() {
	elevio.SetMotorDirection(elevio.MD_Down)
	elevator.dirn = elevio.MD_Down
	elevator.behaviour = EB_Moving
}



func Fsm_onRequestButtonPress(btn_floor int, btn_type elevio.ButtonType) {
	pc, _, _, _ := runtime.Caller(0)
	functionName := runtime.FuncForPC(pc).Name()
	fmt.Printf("\n\n%s(%d, %s)\n", functionName, btn_floor, Elevio_button_toString(btn_type))
	Elevator_print(elevator)

	switch elevator.behaviour {
	case EB_DoorOpen:
		if Requests_shouldClearImmediatly(elevator, btn_floor, btn_type) {
			Timer_start((elevator.config.doorOpenDuration_s))
		} else {
			elevator.request[btn_floor][btn_type] = 1
		}
		break

	case EB_Moving:
		elevator.request[btn_floor][btn_type] = 1
		break
	
	case EB_Idle:
		elevator.request[btn_floor][btn_type] = 1
		var pair DirnBehaviourPair = Request_chooseDirection(elevator)
		elevator.dirn = pair.dirn
		elevator.behaviour = pair.behaviour

		switch pair.behaviour {
		case EB_DoorOpen:
			elevio.SetDoorOpenLamp(true)
			Timer_start(elevator.config.doorOpenDuration_s)
			elevator = Requests_clearAtCurrentFloor(elevator)
			break
		
		case EB_Moving:
			elevio.SetMotorDirection(elevator.dirn)
			break
		
		case EB_Idle:
			break

		}
		break

	}

	SetAllLights(elevator)
	fmt.Printf("\nNew state:\n")
	Elevator_print(elevator)
}

func Fsm_onFloorArrival(newFloor int) {
	
	pc, _, _, _ := runtime.Caller(0)
	functionName := runtime.FuncForPC(pc).Name()
	fmt.Printf("\n\n%s(%d)\n", functionName, newFloor)
	Elevator_print(elevator)
	
	elevator.floor = newFloor
	elevio.SetFloorIndicator(elevator.floor)
	
	switch elevator.behaviour {
	case EB_Moving:
		if Requests_shouldStop(elevator) {
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevio.SetDoorOpenLamp(true)
			elevator = Requests_clearAtCurrentFloor(elevator)
			Timer_start(elevator.config.doorOpenDuration_s)
			SetAllLights(elevator)
			elevator.behaviour = EB_DoorOpen
		}
		break
	default:
		break
	}
	fmt.Printf("\nNew state:\n")
	Elevator_print(elevator)
}


func Fsm_onDoorTimeout() {
	pc, _, _, _ := runtime.Caller(0)
	functionName := runtime.FuncForPC(pc).Name()
	fmt.Printf("\n\n%s()\n", functionName)

	Elevator_print(elevator)

	switch elevator.behaviour {
	case EB_DoorOpen:
		var pair DirnBehaviourPair = Request_chooseDirection(elevator)
		elevator.dirn = pair.dirn
		elevator.behaviour = pair.behaviour

		switch elevator.behaviour {
		case EB_DoorOpen:
			Timer_start(elevator.config.doorOpenDuration_s)
			elevator = Requests_clearAtCurrentFloor(elevator)
			SetAllLights(elevator)
			break
		
		case EB_Moving, EB_Idle:
			elevio.SetDoorOpenLamp(false)
			elevio.SetMotorDirection(elevator.dirn)

			break
		}
		break
	default:
		break
	}
	fmt.Printf("\nNew state:\n")
	Elevator_print(elevator)
}
