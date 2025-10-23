package tools

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

const (
	ARRAY_FIRST_SYMBOL  = "*"
	STRING_FIRST_SYMBOL = "$"
)

const (
	ARRAY_TYPE  = 1
	STRING_TYPE = 2
)

type Parser struct {
	reader *bufio.Reader
}

func New(conn net.Conn) *Parser {
	return &Parser{
		reader: bufio.NewReader(conn),
	}
}

func (p *Parser) Parse() ([]string, error) {
	msg, err := p.readTillEOF()
	if err != nil {
		fmt.Errorf(err.Error())
	}

	if len(msg) == 0 {
		return []string{}, errors.New("Пустое сообщение")
	}

	str := string(msg)
	fmt.Println(str)
	commandType := string(str[0])
	if commandType != ARRAY_FIRST_SYMBOL {
		return []string{}, errors.New("Неверный формат команды")
	}

	elementCount, e := strconv.Atoi(string(str[1]))
	if e != nil {
		return []string{}, nil
	}
	if elementCount == 0 {
		return []string{}, nil
	}

	elements := make([]string, 0, elementCount)

	for i := 0; i < elementCount; i++ {
		element, err := p.readTillEOF()
		if err != nil {
			return []string{}, err
		}
		elType := p.getElementType(element)
		if elType == STRING_TYPE {
			el := p.getString()
			elements = append(elements, el)
		}
	}

	return elements, nil
}

func (p *Parser) readTillEOF() ([]byte, error) {
	msg, err := p.reader.ReadBytes('\n')
	if err != nil {
		fmt.Println(err)
		return []byte{}, err
	}

	return msg, nil
}

func (p *Parser) getElementType(el []byte) uint8 {
	firtsSymbol := string(el[0])

	switch firtsSymbol {
	case ARRAY_FIRST_SYMBOL:
		return ARRAY_TYPE
	case STRING_FIRST_SYMBOL:
		return STRING_TYPE
	default:
		return STRING_TYPE
	}
}

func (p *Parser) getString() string {
	bytes, err := p.readTillEOF()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return strings.TrimSpace(string(bytes))
}
