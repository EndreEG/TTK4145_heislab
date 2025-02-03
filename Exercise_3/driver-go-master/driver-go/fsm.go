package fsm

import (
	"Driver-go/elevio"
	"Driver-go/request"
	"Driver-go/elevator"
	"runtime"
	"Driver-go/timer"
	"fmt"
)

func fsm_init() {
	elevator = elevator_uninitialized()

	
}

func setAllLights(e *Elevator){
	for floor := 0; floor < NumFloors; floor++ {
		for btn = 0; btn < NumButtons; btn++ {
			SetButtonLamp(btn, floor, e.request[floor][btn])
		}
	}
}

func fsm_onInitBetweenFloors(e *Elevator) {
	SetMotorDirection(D_Down)
	e.dirn = D_Down
	e.behaviour = EB_Moving
}

func fsm_onRequestButtonPress(e *Elevator, btn_floor int, btn_type Button) {
	for{
		select {
		case EB_DoorOpen:
			if request_shouldClearImmediatly(e *Elevator, btn_type, btn_floor) == true {
				timer_start()
			} else{
				e.request[btn_floor][btn_type] := 1
			}
			break

		case EB_Moving:
			e.request[btn_floor][btn_type] := 1
			break
		
		case EB_Idle:
			e.request[btn_floor][btn_type] := 1
			DirnBehaviourPair pair := request_chooseDirection(e)
			setMotorDirection(pair.dirn)
			e.behaviour = pair.behavior
			for{
				select{
				case EB_DoorOpen:
					SetDoorOpenLamp(1)
					timer_start(3)
					e = request_clearAtCurrentFloor(e)
					break
				
				case EB_Moving:
					SetMotorDirection(e.dirn)
					break

				case EB_Idle:
					break
				}
				break
			}

			setAllLights(e)

		}
	}
}





// FsmOnFloorArrival handles the logic when the elevator arrives at a new floor.
// It updates the elevator's state, checks if it should stop, and manages lights and timers.
func FsmOnFloorArrival(newFloor int, e *elevator.Elevator) {
	// Print the function name and the new floor for debugging purposes.
	fmt.Printf("\n\nFsmOnFloorArrival(%d)\n", newFloor)

	// Print the current state of the elevator.
	elevator.ElevatorPrint(*e)

	// Update the elevator's current floor.
	e.Floor = newFloor

	// Update the floor indicator light on the hardware.
	elevio.SetFloorIndicator(e.Floor)

	// Handle behavior based on the elevator's current state.
	switch e.Behaviour {
	case elevator.EB_Moving:
		// If the elevator is moving, check if it should stop at the current floor.
		if request.ShouldStop(e) {
			// Stop the motor.
			elevio.SetMotorDirection(elevio.MD_Stop)

			// Turn on the door open lamp.
			elevio.SetDoorOpenLamp(true)

			// Clear requests at the current floor.
			*e = request_clearAtCurrentFloor(e)

			// Start the door open timer.
			timer.Start(e.Config.DoorOpenDuration_s)

			// Update all button lights to reflect the current state of requests.
			setAllLights(e)

			// Change the elevator's behavior to DoorOpen.
			e.Behaviour = elevator.EB_DoorOpen
		}
	default:
		// For other behaviors (e.g., Idle, DoorOpen), do nothing.
	}

	// Print the new state of the elevator after processing.
	fmt.Println("\nNew state:")
	elevator.ElevatorPrint(*e)
}

// setAllLights updates all button lights based on the elevator's request matrix.
func setAllLights(e *elevator.Elevator) {
	// Iterate through all floors and buttons.
	for floor := 0; floor < elevator.NumFloors; floor++ {
		for btn := 0; btn < elevator.NumButtons; btn++ {
			// Set the button lamp to on if there is a request, otherwise turn it off.
			elevio.SetButtonLamp(elevio.ButtonType(btn), floor, e.Requests[floor][btn] == 1)
		}
	}
}