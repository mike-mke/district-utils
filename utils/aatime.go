package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type AaTime struct {
	Hour      string
	Minute    string
	AmPm      string
	IsAm      bool
	NumHour   int
	NumMinute int
	IsValid   bool
}

func NewAaTime(t string) (AaTime, error) {
	time := AaTime{
		Hour:      "23",
		Minute:    "59",
		AmPm:      "PM",
		IsAm:      false,
		NumHour:   23,
		NumMinute: 59,
		IsValid:   false,
	}

	c := strings.Index(t, ":")
	if c == -1 {
		return time, errors.New("invalid time string:" + t)
	}

	time.Hour = t[0:c]

	rest := t[c+1:] // everything after the colon, e.g. "00 AM" or "45"
	s := strings.Index(rest, " ")

	hasAmPm := s != -1
	validAmPm := true
	if hasAmPm {
		// 12-hour format: "07:00 AM", "12:30 PM"
		time.Minute = rest[0:s]
		time.AmPm = strings.TrimSpace(rest[s+1:])
		switch time.AmPm {
		case "AM":
			time.IsAm = true
		case "PM":
			time.IsAm = false
		default:
			fmt.Printf("Unknown AM/PM value: %s \n", time.AmPm)
			validAmPm = false
		}
	} else {
		// 24-hour format: "17:45"
		time.Minute = rest
	}

	numHour, err := strconv.ParseInt(time.Hour, 10, 32)
	if err != nil {
		fmt.Printf("err: %s", err)
		return time, err
	} else {
		time.NumHour = int(numHour)
	}

	numMinute, err2 := strconv.ParseInt(time.Minute, 10, 16)
	if err2 != nil {
		fmt.Printf("err: %s", err2)
		return time, err
	} else {
		time.NumMinute = int(numMinute)
	}

	// Validate before normalization.
	// Minutes are always 0–59 regardless of format.
	// 12-hour format: hour must be 1–12.
	// 24-hour format: hour must be 0–23.
	if err != nil || err2 != nil || !validAmPm {
		return time, errors.New("invalid time string:" + t) // IsValid stays false
	}
	if time.NumMinute < 0 || time.NumMinute > 59 {
		return time, errors.New("Invalid minute value: " + time.Minute)
	}
	if hasAmPm {
		if time.NumHour < 1 || time.NumHour > 12 {
			return time, errors.New("Invalid hour value: " + time.Hour)
		}
	} else {
		if time.NumHour < 0 || time.NumHour > 23 {
			fmt.Printf("Invalid 24-hour clock hour value: %d\n", time.NumHour)
			return time, errors.New("Invalid 24-hour clock hour value: " + time.Hour)
		}
	}

	// Normalize 12-hour to 24-hour after validation.
	if hasAmPm {
		//   12:xx AM → 00:xx  (midnight)
		//   12:xx PM → 12:xx  (noon, already correct)
		//    1–11 PM → add 12
		if time.IsAm && time.NumHour == 12 {
			time.NumHour = 0
		} else if !time.IsAm && time.NumHour < 12 {
			time.NumHour += 12
		}
	} else {
		// Derive AmPm/IsAm from the 24-hour value for display consistency
		if time.NumHour < 12 {
			time.AmPm = "AM"
			time.IsAm = true
		} else {
			time.AmPm = "PM"
			time.IsAm = false
		}
	}

	time.IsValid = true
	return time, nil
}

func (time AaTime) Print() {
	fmt.Printf("Time: %02d:%02d %s (%02d:%02d)\n", time.NumHour, time.NumMinute, time.AmPm, time.NumHour, time.NumMinute)
}

func (time AaTime) GetKey() string {
	return fmt.Sprintf("%02d:%02d", time.NumHour, time.NumMinute)
}
