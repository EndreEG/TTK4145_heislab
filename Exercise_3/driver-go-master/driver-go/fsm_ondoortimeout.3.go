package fsm

import (
	"fmt"
	"runtime"
)

func fsm_onDoorTimeout() {
	pc, _, _, _ := runtime.Caller(0)
	functionName := runtime.FuncForPC(pc).Name()
	fmt.printf("\n\n%s()\n", functionName)

	elevator_print(elevator)

	switch elevator.behaviour {
	case EB_DoorOpen:
		var pair DirnBehaviourPair = request_chooseDirection(elevator)
		elevator.dirn = pair.dirn
		elevator.behaviour = pair.behaviour

		switch elevator.behaviour {
		case EB_DoorOpen:
			timer_start(elevator.config.doorOpenDuration_s)
			elevator = request_clearAtCurrentFloor(elevator)
			setAllLights(elevator)
			break
		case EB_Moving, EB_Idle:
			outputDevice.doorLight(0)
			outputDevice.motorDirection(elevator.dirn)
			break
		}
		break
	default:
		break
	}
	fmt.Printf(printf("\nNew state:\n"))
	elevator_print(elevator)
}
