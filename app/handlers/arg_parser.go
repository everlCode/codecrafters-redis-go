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

func (ap ArgParser) ParseTimeout(timeoutValue string, isSeconds bool) time.Duration {
	v, err := strconv.ParseFloat(timeoutValue, 64)
	if err != nil || v == 0 {
		return 0
	}
	if isSeconds {
		return time.Duration(v * float64(time.Second))
	}

	return time.Duration(v) * time.Millisecond
}
