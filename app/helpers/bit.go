package helpers

func SetBit(storage []byte, offset int, value bool) ([]byte, int, error) {
	const byte_lenght = 8
	byteNumber := ((offset + byte_lenght - 1) / byte_lenght)
	stringLenght := len(storage)

	if stringLenght < byteNumber {
		storage = append(storage, make([]byte, byteNumber - stringLenght)...)
	}
	bitNumber := offset % 8

	byteToChange := storage[byteNumber - 1]
	var originalValue byte
	originalValue = byteToChange & (1 << bitNumber)

	if value {
		byteToChange |= 1 << bitNumber
	} else {
		byteToChange &^= 1 << bitNumber
	}
	storage[byteNumber - 1] = byteToChange

	var originalValueInt int
	if int(originalValue) > 0 {
		originalValueInt = 1
	}
		
	return storage, originalValueInt, nil
}