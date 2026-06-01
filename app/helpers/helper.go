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

func IncrementStreamId(id string) string {
	parts := strings.Split(id, "-")
	if len(parts) < 2 {
		value, _ :=  strconv.Atoi(id)
		value += 1

		return strconv.Itoa(value)
	}

	indetificator, _ := strconv.Atoi(parts[1])
	indetificator += 1
	indetificatorString := strconv.Itoa(indetificator)

	return strings.Join([]string{parts[0], indetificatorString}, "-")
}

func GetDataByStartEnd[T any](data []T, start int, end int) []T {
	lenght := len(data)
	if start < 0 {
		start = lenght + start
		if start < 0 {
			start = 0
		}
	}

	if end < -1 {
		end = lenght + end + 1
	} else if end < lenght && end > 0 {
		end += 1
	}

	if start > end && end >= 0 {
		return []T{}
	}

	if start >= lenght {
		return []T{}
	}

	if end == -1 || end >= lenght || end == 0 {
		data = data[start:]
	} else {
		data = data[start:end]
	}

	return data
}