package helpers

func SetBit(storage []byte, offset int, value bool) ([]byte, int, error) {
	byteIndex := byteNumberByOffset(offset)
	
	stringLenght := len(storage)

	if stringLenght < (byteIndex + 1) || stringLenght == 0 {
		var size int
		if byteIndex == 0 {
			size = 1
		} else {
			size = (byteIndex + 1) - stringLenght
		}
		storage = append(storage, make([]byte, size)...)
	}
	bitNumber := offset % 8

	byteToChange := storage[byteIndex]
	var originalValue bool
	originalValue = isTrueBit(byteToChange, bitNumber)

	if value {
		byteToChange |= 128 >> bitNumber
	} else {
		byteToChange &^= 128 >> bitNumber
	}
	storage[byteIndex] = byteToChange

	var originalValueInt int
	if originalValue {
		originalValueInt = 1
	}
		
	return storage, originalValueInt, nil
}

func GetBit(storage []byte, offset int) (int) {
	byteIndex := byteNumberByOffset(offset)
	stringLenght := len(storage)

	if stringLenght < (byteIndex + 1) {
		return 0
	}
	bitNumber := offset % 8

	byteToChange := storage[byteIndex]
	var originalValue bool
	originalValue = isTrueBit(byteToChange, bitNumber)

	var originalValueInt int
	if originalValue {
		originalValueInt = 1
	}
		
	return originalValueInt
}

func CountTrueBit(storage []byte) int {
	var counter int

	for _, curByte := range storage {
		for i := range 8 {
			if isTrueBit(curByte, i) {
				counter++
			}
		}
	}
	
	return counter
}

func isTrueBit(b byte, offset int) bool {
	if offset < 0 || offset > 7 {
		return false
	}

	val := b & (128 >> offset)

	return val > 0
}

func byteNumberByOffset(offset int) int {
	var byteNumber int
	if offset <= 7 {
		byteNumber = 0
	} else {
		byteNumber = (offset / 8)
	}

	return byteNumber
}

func BitAndOp(a []byte, b []byte) []byte {
	var res []byte

	for i, curByte := range a {
		if i + 1 > len(b) {
			break;
		}
		maskResult := curByte & b[i]
		res = append(res, maskResult)
	}

	return res
}