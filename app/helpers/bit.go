package helpers

const byte_lenght = 8

func SetBit(storage []byte, offset int, value bool) ([]byte, int, error) {
	byteNumber := ((offset + byte_lenght - 1) / byte_lenght)
	stringLenght := len(storage)

	if stringLenght < byteNumber || stringLenght == 0 {
		var size int
		if byteNumber == 0 {
			size = 1
		} else {
			size = byteNumber - stringLenght
		}
		storage = append(storage, make([]byte, size)...)
	}
	bitNumber := offset % 8

	if byteNumber == 0 {
		byteNumber++
	}
	byteToChange := storage[byteNumber - 1]
	var originalValue byte
	originalValue = byteToChange & (128 >> bitNumber)

	if value {
		byteToChange |= 128 >> bitNumber
	} else {
		byteToChange &^= 128 >> bitNumber
	}
	storage[byteNumber - 1] = byteToChange

	var originalValueInt int
	if int(originalValue) > 0 {
		originalValueInt = 1
	}
		
	return storage, originalValueInt, nil
}

func GetBit(storage []byte, offset int) (int) {
	byteNumber := ((offset + byte_lenght - 1) / byte_lenght)
	if byteNumber == 0 {
		byteNumber++
	}
	stringLenght := len(storage)

	if stringLenght < byteNumber {
		return 0
	}
	bitNumber := offset % 8

	byteToChange := storage[byteNumber - 1]
	var originalValue byte
	originalValue = byteToChange & (128 >> bitNumber)

	var originalValueInt int
	if int(originalValue) > 0 {
		originalValueInt = 1
	}
		
	return originalValueInt
}