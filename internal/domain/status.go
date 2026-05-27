package domain

import "strings"

type Status string

const (
	StatusOn  Status = "on"
	StatusOff Status = "off"
)

func ParseStatus(s string) (Status, bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "on":
		return StatusOn, true
	case "off":
		return StatusOff, true
	default:
		return "", false
	}
}

func (s Status) IsOn() bool {
	return s == StatusOn
}
