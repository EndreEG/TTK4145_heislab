package main

import (
	"Heis/elevio"
	"fmt"
	"runtime"
)

var elevator Elevator

func setAllLights(es Elevator) {
	bu := []elevio.ButtonType{elevio.BT_HallUp, elevio.BT_HallDown, elevio.BT_Cab}
	for floor := 0; floor < NumFloors; floor++ {
		for _, B := range bu {
			elevio.SetButtonLamp(B, floor, es.request[floor][B] == 1)
		}
	}
}

func fsm_onInitBetweenFloors() {
	elevio.SetMotorDirection(elevio.MD_Down)
	elevator.dirn = elevio.MD_Down
	elevator.behaviour = EB_Moving
}

func fsm_onRequestButtonPress(btn_floor int, btn_type elevio.ButtonType) {
	pc, _, _, _ := runtime.Caller(0)
	functionName := runtime.FuncForPC(pc).Name()
	fmt.Printf("\n\n%s(%d, %s)\n", functionName, btn_floor, elevio_button_toString(btn_type))
	elevator_print(elevator)

	switch elevator.behaviour {
	case EB_DoorOpen:
		if requests_shouldClearImmediatly(elevator, btn_floor, btn_type) {
			timer_start((elevator.config.doorOpenDuration_s))
		} else {
			elevator.request[btn_floor][btn_type] = 1
		}
		break
	case EB_Moving:
		elevator.request[btn_floor][btn_type] = 1
		break
	case EB_Idle:
		elevator.request[btn_floor][btn_type] = 1
		var pair DirnBehaviourPair = request_chooseDirection(elevator)
		elevator.dirn = pair.dirn
		elevator.behaviour = pair.behaviour

		switch pair.behaviour {
		case EB_DoorOpen:
			elevio.SetDoorOpenLamp(true)
			timer_start(elevator.config.doorOpenDuration_s)
			elevator = requests_clearAtCurrentFloor(elevator)
			break
		case EB_Moving:
			elevio.SetMotorDirection(elevator.dirn)
			break
		case EB_Idle:
			break

		}
		break

	}
	setAllLights(elevator)
	fmt.Printf("\nNew state:\n")
	elevator_print(elevator)
}

func fsm_onFloorArrival(newFloor int) {
	pc, _, _, _ := runtime.Caller(0)
	functionName := runtime.FuncForPC(pc).Name()
	fmt.Printf("\n\n%s(%d)\n", functionName, newFloor)
	elevator_print(elevator)
	elevator.floor = newFloor
	elevio.SetFloorIndicator(elevator.floor)
	switch elevator.behaviour {
	case EB_Moving:
		if requests_shouldStop(elevator) {
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevio.SetDoorOpenLamp(true)
			elevator = requests_clearAtCurrentFloor(elevator)
			timer_start(elevator.config.doorOpenDuration_s)
			setAllLights(elevator)
			elevator.behaviour = EB_DoorOpen
		}
		break
	default:
		break
	}
	fmt.Printf("\nNew state:\n")
	elevator_print(elevator)
}

func fsm_onDoorTimeout() {
	pc, _, _, _ := runtime.Caller(0)
	functionName := runtime.FuncForPC(pc).Name()
	fmt.Printf("\n\n%s()\n", functionName)

	elevator_print(elevator)

	switch elevator.behaviour {
	case EB_DoorOpen:
		var pair DirnBehaviourPair = request_chooseDirection(elevator)
		elevator.dirn = pair.dirn
		elevator.behaviour = pair.behaviour

		switch elevator.behaviour {
		case EB_DoorOpen:
			timer_start(elevator.config.doorOpenDuration_s)
			elevator = requests_clearAtCurrentFloor(elevator)
			setAllLights(elevator)
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
	elevator_print(elevator)
}
