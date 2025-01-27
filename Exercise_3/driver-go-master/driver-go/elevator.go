//Elevator.h
type ElevatorBehaviour int
const (EB_Idle = iota
		EB_DoorOpen
		EB_Moving)

type clearRequestVariant int
const (CV_All = iota
		CV_InDirn)


type Elevator struct{
	floor int
	dirn Dirn
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

func elevator_print(es Elevator) {
	fmt.Printf("  +--------------------+\n")
	fmt.Printf("|   floor = %s   |\n    |dirn = %s     |\n    behav = %n", es.floor,
		elevio_dirn_toString(es.dirn), eb_toString(es.behaviour))
	fmt.printf("  +--------------------+\n")
	fmt.printf("  |  | up  | dn  | cab |\n")
	for (f := N_FLOORS-1;f>=0;f--){
		fmt.printf
	}
}





func elevator_uninitialized()Elevator{
	el := Elevator{floor : -1,
					dirn : D_Stop,
					behavior : EB_Idle,
					config : {clearRequestVariant = CV_All ,doorOpenDuration_s = 3.0}
				}
}
