package main

import (
	"fmt"
	"testing"
)

func TestRandomFloat32(t *testing.T) {
	value := randomFloat32()
	result := fmt.Sprintf("%T", value)
	want := "float32"

	if result != want {
		t.Fatalf("randomFloat32() does not return float32")
	}

}

func TestRandomIntMinMax(t *testing.T) {
	value := randomIntMinMax(10, 20)

	if value < 10 || value > 20 {
		t.Fatalf("randomIntMinMax return out of bounds value")
	}

}

func TestMask(t *testing.T) {
	result := mask()
	want := "********"

	if result != want {
		t.Fatalf("mask() does not return %s but %s", want, result)
	}
}
