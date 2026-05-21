package handlers

import (
	"strconv"
	"time"
)

const (
	PX = "px"
)

type ArgParser struct {
}

func New() *ArgParser {
	return &ArgParser{}
}

func (ap ArgParser) ParseTimeout(timeoutValue string) time.Duration {
	var timeout time.Duration = 0

	v, err := strconv.ParseFloat(timeoutValue, 64)
	if err != nil || v == 0 {
		return 0
	}

	timeout = time.Duration(v * float64(time.Second))

	return timeout
}
