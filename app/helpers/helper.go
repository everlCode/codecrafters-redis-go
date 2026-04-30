package helpers

import (
	"strconv"
	"strings"
)

func GetStreamIdParts(id string) (int, int) {
	parts := strings.Split(id, "-")
	if len(parts) < 2 {
		return 0, 0
	}

	miliseconds, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0
	}

	indetificator, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0
	}

	return miliseconds, indetificator
}
