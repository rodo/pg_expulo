package main

import (
	"testing"
	"fmt"
)


func TestRandomFloat32(t *testing.T) {

	value := randomFloat32()
	result := fmt.Sprintf("%T", value)
	want := "float32"


	if result != want {
		t.Fatalf("randomFloat32() does not return float32")
	}

}
