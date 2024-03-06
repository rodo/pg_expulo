package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
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

func TestAllowedGenerators(t *testing.T) {
	result := allowedGenerators()

	// Assert that result is of type []string
	assert.True(t, slices.Contains(result, "serial"), "Expected result of type string")
}
