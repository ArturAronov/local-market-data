package utils

import (
	"fmt"
	"strconv"
	"time"
)

type Date struct {
	Year  string
	Month string
	Day   string
}

func DateParser(date Date) (time.Time, error) {
	if len(date.Year) != 4 {
		return time.Time{}, fmt.Errorf("Wrong year format")
	}

	if len(date.Month) != 2 {
		return time.Time{}, fmt.Errorf("Wrong month format")
	}

	if len(date.Day) != 2 {
		return time.Time{}, fmt.Errorf("Wrong day format")
	}

	y, yErr := strconv.Atoi(date.Year)
	if yErr != nil || y < 1900 || y > 2100 {
		return time.Time{}, fmt.Errorf("Invalid year: %w", yErr)
	}

	m, mErr := strconv.Atoi(date.Month)
	if mErr != nil || m < 1 || m > 12 {
		return time.Time{}, fmt.Errorf("Invalid month: %w", yErr)
	}

	d, dErr := strconv.Atoi(date.Month)
	if dErr != nil || d < 1 || d > 31 {
		return time.Time{}, fmt.Errorf("Invalid day: %w", yErr)
	}

	t := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
	if t.Year() != y || t.Month() != time.Month(m) || t.Day() != d {
		return time.Time{}, fmt.Errorf(
			"Invalid calendar date: %s-%s-%s",
			date.Year,
			date.Month,
			date.Day,
		)
	}

	return t, nil
}
