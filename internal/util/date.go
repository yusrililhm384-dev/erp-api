package util

import (
	"strings"
	"time"
)

type JSONDate struct {
	time.Time
}

func (jd *JSONDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")

	if s == "null" || s == "" {
		return nil
	}

	t, err := time.Parse(time.DateOnly, s)

	if err != nil {
		return err
	}

	jd.Time = t

	return nil
}

type JSONNullDate struct {
	*time.Time
}

func (jnd *JSONNullDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")

	if s == "null" || s == "" {
		jnd.Time = nil
		return nil
	}

	t, err := time.Parse(time.DateOnly, s)

	if err != nil {
		return err
	}

	jnd.Time = &t

	return nil
}
