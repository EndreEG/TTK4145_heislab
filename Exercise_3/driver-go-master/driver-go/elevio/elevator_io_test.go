package elevio_test

import (
	"Driver-go/elevio"
	"testing"
)

func TestGetFloor(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := elevio.GetFloor()
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("GetFloor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInit(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		addr      string
		numFloors int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			elevio.Init(tt.addr, tt.numFloors)
		})
	}
}

func TestInit(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		addr      string
		numFloors int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			elevio.Init(tt.addr, tt.numFloors)
		})
	}
}
