package main

import (
	"Heis/elevio"
	"fmt"
	"time"
)

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
        Fsm_onInitBetweenFloors();
    }

	for {

		select {
		case a := <-drv_buttons:
			fmt.Printf("%+v\n", a)
			Fsm_onRequestButtonPress(a.Floor, a.Button)

		case a := <-drv_floors:
			fmt.Printf("%+v\n", a)
			Fsm_onFloorArrival(a)

		case a := <- timeout:
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
			fmt.Printf("%+v\n", a)
			elevio.SetMotorDirection(elevio.MD_Stop)
			for f := 0; f < numFloors; f++ {
				for b := elevio.ButtonType(0); b < 3; b++ {
					elevio.SetButtonLamp(b, f, false)
				}
			}
		}

		time.Sleep((500 * time.Duration(inputPollRate)))
}
