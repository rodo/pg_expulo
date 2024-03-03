package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// toolPat() without one SQL function at the end
// 1 row with 2 columns
func TestToolPatFuncOneTwo(t *testing.T) {

	result := toolPat(1, []string{"$1", "now()"})
	want := "($1,now())"

	if result != want {
		t.Fatalf("Return %s in place of %s", result, want)
	}
}

// toolPat() without SQL function
// 1 row with 1 columns
func TestToolPatOneOne(t *testing.T) {

	result := toolPat(1, []string{"$1"})
	want := "($1)"

	if result != want {
		t.Fatalf("Return %s in place of %s", result, want)
	}
}

// toolPat() without SQL function
// 2 rows with 2 columns
func TestToolPat(t *testing.T) {

	result := toolPat(2, []string{"$1", "$2"})
	want := "($1,$2),($3,$4)"

	if result != want {
		t.Fatalf("Return %s in place of %s", result, want)
	}
}

// toolPat() without SQL function
// 1 row with 3 columns
func TestToolPatSix(t *testing.T) {

	result := toolPat(1, []string{"$1", "$2", "$3"})
	want := "($1,$2,$3)"

	if result != want {
		t.Fatalf("Return %s in place of %s", result, want)
	}
}

func TestToolPatNine(t *testing.T) {

	result := toolPat(3, []string{"$1", "$2", "$3"})
	want := "($1,$2,$3),($4,$5,$6),($7,$8,$9)"

	if result != want {
		t.Fatalf("Return %s in place of %s", result, want)
	}
}

func TestToolPatFuncTwo(t *testing.T) {

	result := toolPat(2, []string{"$1", "now()"})
	want := "($1,now()),($2,now())"

	if result != want {
		t.Fatalf("Return %s in place of %s", result, want)
	}
}

// toolPat() without one SQL function in the middle
// 2 rows with 4 columns
func TestToolPatFuncTwoFour(t *testing.T) {

	result := toolPat(2, []string{"$1", "now()", "$2", "$3"})
	want := "($1,now(),$2,$3),($4,now(),$5,$6)"

	if result != want {
		t.Fatalf("Return %s in place of %s", result, want)
	}
}

func TestToolPatFunc(t *testing.T) {

	result := toolPat(2, []string{"$1", "now()", "$2"})
	want := "($1,now(),$2),($3,now(),$4)"

	if result != want {
		t.Fatalf("Return %s in place of %s", result, want)
	}
}

// toolPat() without SQL function
// 3 rows with 2 columns
func TestToolPatNone(t *testing.T) {

	result := toolPat(3, []string{"$1", "$2"})
	want := "($1,$2),($3,$4),($5,$6)"

	if result != want {
		t.Fatalf("Return %s in place of %s", result, want)
	}
}

// toolPat() with only SQL function
// 3 rows with 2 columns
func TestToolPatOnlyFunction(t *testing.T) {

	result := toolPat(3, []string{"now()", "getRandomString()"})
	want := "(now(),getRandomString()),(now(),getRandomString()),(now(),getRandomString())"

	if result != want {
		t.Fatalf("Return %s in place of %s", result, want)
	}
}

// getEnvVariable returns an well set variable
func TestGetEnvVariableSuccess(t *testing.T) {
	// Set up the environment variable for testing
	expectedValue := "test_value"
	envVarName := "TEST_VAR"
	err := os.Setenv(envVarName, expectedValue)
	assert.NoError(t, err, "Error setting up test environment")

	// Call the function being tested
	result := getEnvVariable(envVarName)

	// Assert that the result matches the expected value
	assert.Equal(t, expectedValue, result, "Result does not match expected value")

	// Clean up the environment variable after the test
	err = os.Unsetenv(envVarName)
	assert.NoError(t, err, "Error cleaning up test environment")

}

// getEnvVariable returns an well set variable
func TestReadEnvSuccess(t *testing.T) {
	// Set up the environment variable for testing
	expectedValue := ConnectionParameters{"host", "port", "user", "pass", "db"}
	var err error

	err = os.Setenv("PGFOOHOST", "host")
	assert.NoError(t, err, "Error setting up test environment")
	err = os.Setenv("PGFOOPORT", "port")
	assert.NoError(t, err, "Error setting up test environment")
	err = os.Setenv("PGFOOUSER", "user")
	assert.NoError(t, err, "Error setting up test environment")
	err = os.Setenv("PGFOOPASSWORD", "pass")
	assert.NoError(t, err, "Error setting up test environment")
	err = os.Setenv("PGFOODATABASE", "db")
	assert.NoError(t, err, "Error setting up test environment")

	// Call the function being tested
	result := readEnv("FOO")

	// Assert that the result matches the expected value
	assert.Equal(t, expectedValue, result, "Result does not match expected value")

	// Clean up the environment variable after the test
	err = os.Unsetenv("PGFOOHOST")
	assert.NoError(t, err, "Error setting up test environment")
	err = os.Unsetenv("PGFOOPORT")
	assert.NoError(t, err, "Error setting up test environment")
	err = os.Unsetenv("PGFOOUSER")
	assert.NoError(t, err, "Error setting up test environment")
	err = os.Unsetenv("PGFOOPASSWORD")
	assert.NoError(t, err, "Error setting up test environment")
	err = os.Unsetenv("PGFOODATABASE")
	assert.NoError(t, err, "Error setting up test environment")

	assert.NoError(t, err, "Error cleaning up test environment")
}
