package main

import (
	"Heis/elevio"
	"Heis/network"
	"fmt"
	"os"
	"time"
)

type State struct {
	Elevator_id      int
	Elevator_floor   int
	Elevator_dir     int
	ElevatorBehaviour int
	Elevator_request [NumFloors][NumButtons]int
}

var state State

func main() {
	numFloors := 4

	elevio.Init("localhost:15657", numFloors)

	//var d elevio.MotorDirection = elevio.MD_Up
	//elevio.SetMotorDirection(d)

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
	timeout := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)
	go PollTimer(timeout)

	inputPollRate := 25

	if elevio.GetFloor() == -1 {
		Fsm_onInitBetweenFloors()
	}

	for {

		select {
		case a := <-drv_buttons:
			fmt.Printf("%+v\n", a)
			Fsm_onRequestButtonPress(a.Floor, a.Button)

		case a := <-drv_floors:
			fmt.Printf("%+v\n", a)
			Fsm_onFloorArrival(a)

		case a := <-timeout:
			fmt.Printf("%+v\n", a)
			Timer_stop()
			Fsm_onDoorTimeout()

		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			if a && elevator.behaviour == EB_DoorOpen {
				elevio.SetMotorDirection(elevio.MD_Stop)
				for {
					b := <-drv_obstr
					if !b {
						break
					}
				}
			}

		case a := <-drv_stop:
			elevio.SetStopLamp(true)
			fmt.Printf("%+v\n", a)
			elevio.SetMotorDirection(elevio.MD_Stop)
			for floor := 0; floor < numFloors; floor++ {
				for btn := elevio.ButtonType(0); btn < 3; btn++ {
					elevio.SetButtonLamp(btn, floor, false)
					elevator.request[floor][btn] = 0
				}
			}
			elevio.SetStopLamp(false)
			os.Exit(0)
		}

		state.Elevator_id = 1
		state.Elevator_dir = int(elevator.dirn)
		state.ElevatorBehaviour = int(elevator.behaviour)
		state.Elevator_floor = elevator.floor
		state.Elevator_request = elevator.request

		network.SendElevatorState(state)

		time.Sleep((500 * time.Duration(inputPollRate)))
	}
}
