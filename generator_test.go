package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomString(t *testing.T) {
	value := randomString()
	assert.IsType(t, "string", value, "Expected result of type int32")
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

func TestRandomInt(t *testing.T) {
	result := randomInt()
	assert.IsType(t, int32(0), result, "Expected result of type int32")
}

func TestRandomInt32(t *testing.T) {
	result := randomInt32()
	assert.IsType(t, int32(0), result, "Expected result of type int32")
}

func TestRandomInt64(t *testing.T) {
	result := randomInt64()
	assert.IsType(t, int64(0), result, "Expected result of type int64")
}

func TestRandomFloat(t *testing.T) {
	result := randomFloat()
	assert.IsType(t, float32(0), result, "Expected result of type float32")
}

func TestRandomFloat32(t *testing.T) {
	result := randomFloat32()
	assert.IsType(t, float32(0), result, "Expected result of type float32")
}

func TestRandomFloat64(t *testing.T) {
	result := randomFloat64()
	assert.IsType(t, float64(0), result, "Expected result of type float64")
}

func TestFakeEmail(t *testing.T) {
	result := genericFake{}.FakeEmail()
	assert.NotNil(t, result, "Expected non-nil result")

	// Assert that result is of type string
	assert.IsType(t, "", result, "Expected result of type string")
}

func TestFakeName(t *testing.T) {
	result := genericFake{}.FakeName()
	assert.NotNil(t, result, "Expected non-nil result")

	// Assert that result is of type string
	assert.IsType(t, "", result, "Expected result of type string")
}

func TestFakeFirstName(t *testing.T) {
	result := genericFake{}.FakeFirstName()

	assert.NotNil(t, result, "Expected non-nil result")

	// Assert that result is of type string
	assert.IsType(t, "", result, "Expected result of type string")
}
