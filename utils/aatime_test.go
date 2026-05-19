package utils

import (
	"fmt"
	"testing"
)

func TestTime(t *testing.T) {
	validSamples := []string{"12:00 AM", "12:15 AM", "07:30 AM", "8:00 AM", "12:05 PM", "3:00 PM", "19:00", "00:00", "23:59"}
	for _, t := range validSamples {
		testTimeString(t)
	}

	invalidSamples := []string{"12:79 AM", "32:15 AM", "15:00 AM"}
	for _, t := range invalidSamples {
		testTimeString(t)
	}

	//result := Add(2, 3)
	//expected := 5
	//
	//if result != expected {
	//	t.Errorf("Add(2, 3) = %d; want %d", result, expected)
	//}
}

func testTimeString(t string) {
	aTime, err := NewAaTime(t)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	if aTime.IsValid == false {
		fmt.Println("Invalid time", t)
	} else {
		fmt.Print("Valid time: ")
		aTime.Print()
	}
}
