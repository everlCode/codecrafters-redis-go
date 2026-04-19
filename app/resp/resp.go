package resp

import (
	"bufio"
	"errors"
	"net"
	"strconv"
	"unicode/utf8"
)

const (
	ARRAY  = "*"
	STRING = "+"
	INTEGER = ":"
	BULK   = "$"
	ERROR  = "-"
	CRLF   = "\r\n"
)

type Value struct {
	Type    string
	Integer int
	String  string
	Bulk    string
	Array   []Value
	Expires int64
}

type Parser struct {
	reader *bufio.Reader
}

func New(conn net.Conn) *Parser {
	return &Parser{
		reader: bufio.NewReader(conn),
	}
}

func (p *Parser) Read() (Value, error) {
	_type, err := p.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch string(_type) {
	case ARRAY:
		return p.ReadArray()
	case BULK:
		return p.ReadBulk()
	default:
		return Value{}, errors.New("Неизвестный тип")
	}
}

func (p *Parser) ReadArray() (Value, error) {
	v := Value{
		Type: ARRAY,
	}
	len, err := p.ReadInteger()
	if err != nil {
		return Value{}, err
	}

	elements := make([]Value, 0, len)
	for i := 0; i < len; i++ {
		str, err := p.Read()
		if err != nil {
			return Value{}, err
		}
		elements = append(elements, str)
	}
	v.Array = elements

	return v, nil
}

func (p *Parser) ReadInteger() (int, error) {
	bytes, _, _ := p.reader.ReadLine()
	len, err := strconv.Atoi(string(bytes))
	if err != nil {
		return 0, err
	}

	return len, nil
}

func (p *Parser) ReadBulk() (Value, error) {
	capacity, err := p.ReadInteger()
	if err != nil {
		return Value{}, err
	}

	bytes, _, _ := p.reader.ReadLine()

	str := string(bytes)
	if utf8.RuneCountInString(str) != capacity {
		return Value{}, errors.New("Некорректный запрос! неверное кол-во символов")
	}

	v := Value{
		Type: BULK,
		Bulk: str,
	}

	return v, nil
}

func (v Value) Marshal() []byte {
	switch v.Type {
	case ARRAY:
		var count int
		if v.Array == nil {
			count = -1
		} else {
			count = len(v.Array)
		}
		
		response := []byte(ARRAY + strconv.Itoa(count) + CRLF)
		for _, item := range v.Array {
			response = append(response, item.Marshal()...)
		}

		return response
	case STRING:
		return []byte(STRING + v.String + CRLF)
	case BULK:
		if v.Bulk == "" {
			return []byte(BULK + "-1" + CRLF)
		}
		capacity := len(v.Bulk)

		return []byte(BULK + strconv.Itoa(capacity) + CRLF + v.Bulk + CRLF)
	case INTEGER:
		return []byte(INTEGER + strconv.Itoa(v.Integer) + CRLF)
	case ERROR:
		return []byte(ERROR + v.String + CRLF)
	}

	return []byte{}
}

func Error(msg string) Value {
	return Value{
		Type:   ERROR,
		String: msg,
	}
}
