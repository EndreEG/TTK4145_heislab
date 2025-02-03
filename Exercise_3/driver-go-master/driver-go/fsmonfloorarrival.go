package fsm

import (
	"Driver-go/elevio"
	"fmt"
)

func FsmOnFloorArrival(newFloor int, e *elevator.Elevator) {
	fmt.Printf("\n\nFsmOnFloorArrival(%d)\n", newFloor)
	elevator.ElevatorPrint(*e)

	e.Floor = newFloor

	elevio.SetFloorIndicator(e.Floor)

	switch e.Behaviour {
	case elevator.EB_Moving:
		if request.ShouldStop(e) {
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevio.SetDoorOpenLamp(true)

			*e = request.ClearAtCurrentFloor(e)
			timer.Start(e.Config.DoorOpenDuration_s)
			setAllLights(e)
			e.Behaviour = elevator.EB_DoorOpen
		}
	default:
		// Do nothing for other behaviors
	}

	fmt.Println("\nNew state:")
	elevator.ElevatorPrint(*e)
}

func setAllLights(e *elevator.Elevator) {
	for floor := 0; floor < elevator.NumFloors; floor++ {
		for btn := 0; btn < elevator.NumButtons; btn++ {
			elevio.SetButtonLamp(elevio.ButtonType(btn), floor, e.Requests[floor][btn] == 1)
		}
	}
}
