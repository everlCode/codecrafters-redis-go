package resp

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	ARRAY   = "*"
	STRING  = "+"
	INTEGER = ":"
	BULK    = "$"
	ERROR   = "-"
	CRLF    = "\r\n"
)

type Value struct {
	Type    string
	Integer int
	String  string
	Bulk    string
	Array   []Value
	Bytes   []byte
}

type Parser struct {
	reader *bufio.Reader
	offset int
}

func New(conn net.Conn) *Parser {
	return &Parser{
		reader: bufio.NewReader(conn),
		offset: 0,
	}
}

func (p *Parser) SetOffset(v int) {
	p.offset = v
}

func (p *Parser) AddOffset(v int) {
	p.offset += v
}

func (p *Parser) GetOffset() int {
	return p.offset
}

func (v Value) IsOk() bool {
	if v.Type != STRING {
		return false
	}

	return strings.ToLower(v.String) == "ok"
}

func (p *Parser) Read() (Value, error) {
	_type, err := p.reader.ReadByte()
	p.AddOffset(1)
	if err != nil {
		return Value{}, err
	}

	switch string(_type) {
	case ARRAY:
		return p.ReadArray()
	case BULK:
		return p.ReadBulk()
	case STRING:
		return p.ReadString()
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
	p.AddOffset(len(bytes) + 2)
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
	p.AddOffset(len(bytes) + 2)

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

func (p *Parser) ReadString() (Value, error) {
	bytes, _, _ := p.reader.ReadLine()
	p.AddOffset(len(bytes))

	str := string(bytes)

	v := Value{
		Type:   STRING,
		String: str,
	}

	return v, nil
}

func (p *Parser) ReadRDB() {
	bytes, _, err := p.reader.ReadLine()
	if err != nil {
		panic(err)
	}
	str := string(bytes)
	str = strings.Replace(str, "$", "", 1)
	bodyLength, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}

	var body []byte
	for i := 0; i < bodyLength; i++ {
		b, err := p.reader.ReadByte()
		if err != nil {
			panic(err)
		}
		body = append(body, b)
	}
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

		return []byte(
			BULK +
				strconv.Itoa(capacity) +
				CRLF +
				v.Bulk +
				CRLF,
		)
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

func SimpleString(value string) Value {
	return Value{Type: STRING, String: value}
}

func Bulk(value string) Value {
	return Value{Type: BULK, Bulk: value}
}

func NullBulk() Value {
	return Value{Type: BULK, Bulk: ""}
}

func Integer(value int) Value {
	return Value{Type: INTEGER, Integer: value}
}

func Array(value []any) Value {
	var data []Value = []Value{}

	if len(value) == 0 {
		return EmptyArray()
	}
	for i := 0; i < len(value); i++ {
		switch v := value[i].(type) {
		case []any:
			items, _ := value[i].([]any)
			data = append(data, Array(items))
		case string:
			data = append(data, Value{Type: BULK, Bulk: v})
		case []byte:
			data = append(data, Value{Type: BULK, Bulk: string(v)})
		case int:
			data = append(data, Value{Type: INTEGER, Integer: v})
		case Value:
			data = append(data, v)
		default:
			data = append(data, Value{Type: BULK, Bulk: fmt.Sprintf("%v", v)})
		}
	}

	return Value{Type: ARRAY, Array: data}
}

func ArrayString(value []string) Value {
	var data []Value = []Value{}
	for i := 0; i < len(value); i++ {
		data = append(data, Value{Type: BULK, Bulk: value[i]})
	}
	return Value{Type: ARRAY, Array: data}
}

func EmptyArray() Value {
	return Value{
		Type:  ARRAY,
		Array: []Value{},
	}
}

func NilArray() Value {
	return Value{
		Type:  ARRAY,
		Array: nil,
	}
}

func File(data []byte) Value {
	return Value{
		Type:  BULK,
		Bytes: data,
	}
}

func ParseSlice(input []Value) []string {
	var data []string

	for i := 0; i < len(input); i++ {
		data = append(data, input[i].Bulk)
	}

	return data
}
