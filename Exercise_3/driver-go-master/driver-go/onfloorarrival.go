package fsm

import (
	"Driver-go/elevio"
	"fmt"
)

func OnFloorArrival(newFloor int) {
	fmt.Printf("\n\nOnFloorArrival(%d)\n", newFloor)

	// Update elevator's current floor
	elevator.floor = newFloor

	// Update floor indicator light
	elevio.SetFloorIndicator(elevator.floor)

	switch elevator.Behavior {
	case elevator.Moving:
		// Check if the elevator should stop at this floor
		if requests.ShouldStop(elevator) {
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevio.SetDoorOpenLamp(true)

			// Clear requests at this floor
			elevator = requests.ClearAtCurrentFloor(elevator)

			// Start door open timer
			timer.Start(elevator.Config.DoorOpenDuration)

			// Update button lights
			UpdateAllLights()

			// Change behavior to DoorOpen
			elevator.Behavior = elevator.DoorOpen
		}
	}

	fmt.Println("\nNew state:")
	elevator.PrintState()
}
