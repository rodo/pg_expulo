package main


import (
	"crypto/md5"
	"math/rand"
	"fmt"
	"time"
	"github.com/go-faker/faker/v4"

)

func randomInt() int32 {
	return rand.Int31()
}

func randomIntMinMax(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func randomInt32() int32 {
	return rand.Int31n(100)
}

func randomFloat() float32 {
	return rand.Float32()
}

func randomFloat32() float32 {
	return rand.Float32()
}

func randomFloat64() float64 {
	return rand.Float64()
}

/*
 * String functions
 *
 */
type SomeStruct struct {
	String   string
}

func randomString() string {
	a := SomeStruct{}
	err := faker.FakeData(&a)
	if err != nil {
		fmt.Println(err)
	}
	return a.String
}

func mask() string {
	return "********"
}

func md5signature(String string) string {
	data := []byte(String)
	return fmt.Sprintf("%x", md5.Sum(data))
}

/*
 * Time functions
 *
 */
func randomTimeTZ(timezone string) time.Time {

	location, err := time.LoadLocation(timezone)
	if err != nil {
		fmt.Println("Error:", err)
		loc, _ := time.LoadLocation("UTC")
		return time.Date(1970, 1, 1, 0, 0, 0, 0, loc)
	}

	min := time.Date(1973, 1, 1, 0, 0, 0, 0, location).Unix()
	max := time.Date(2024, 1, 1, 0, 0, 0, 0, location).Unix()
	randomUnix := rand.Int63n(max-min) + min
	return time.Unix(randomUnix, 0)
}
