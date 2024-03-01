package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
