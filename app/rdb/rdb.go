package rdb

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/config"
)

type RDB struct {
	Reader *bufio.Reader
}

type Value struct {
	Type int
	Key string
	Value any
	Expire int64
}

type EncodedValue struct {
    Value string
}

const (
	STRING = 0
)

func New(config *config.Config) *RDB {
	dir := config.Dir
	dbfilename := config.Dbfilename

	if dbfilename == "" {
		return nil
	}
	fmt.Println(os.Getwd())
	fmt.Println(dir + "/" + dbfilename)

	file, err := os.Open(dir + "/" + dbfilename)
	if err != nil {
		return nil
	}

	reader := bufio.NewReader(file)

	return &RDB{
		Reader: reader,
	}
}

func (rdb *RDB) Read() []Value {
	header := rdb.readHeader()
	fmt.Println(header)
	meta := rdb.readMetadata()
	fmt.Println(meta)
	result := rdb.readData()
	rdb.readEnd()

	return result
}

func (rdb *RDB) readHeader() string {
	buf := rdb.readBytes(9)

	return string(buf)
}

func (rdb *RDB) readMetadata() map[string]string {
	metaData := make(map[string]string)
	for rdb.peek() == 0XFA {
		rdb.readByte()
		param := rdb.readEncodedValue()
		fmt.Println(param)
		value := rdb.readEncodedValue()

		metaData[param.Value] = value.Value
	}

	return metaData
}

func (rdb *RDB) readData() []Value {
	b := rdb.readByte()
	if b != 0xFE {
		panic("Invalid rba body!")
	}
	index, _ := rdb.readLenght()
	
	b = rdb.readByte()
	isHashTableSizeInformatiton := b == 0xFB
	if !isHashTableSizeInformatiton {
		panic("Invalid rba body!")
	}
	hashTableValuesSize, _ := rdb.readLenght()
	hashTableExpireSize, _ := rdb.readLenght()

	var result []Value
	for rdb.peek() != 0XFF {
		result = append(result, rdb.readDatabaseValue())
	}
	
	fmt.Println(hashTableValuesSize)
	fmt.Println(hashTableExpireSize)
	fmt.Println(index)

	return result
}

func (rdb *RDB) readEnd() uint64 {
	b := rdb.readByte()
	if b != 0xFF {
		panic("Invalid rba end!")
	}
	buf := rdb.readBytes(8)

	return binary.LittleEndian.Uint64(buf)
}

func (rdb *RDB) readDatabaseValue() Value {
	var expire int64 = 0
	b := rdb.peek()
	if b == 0XFC || b == 0XFD {
		b := rdb.readByte()
		expire =  rdb.readTimestamp(b)
	}

	_type, err := rdb.Reader.ReadByte()
	if err != nil {
		panic(err)
	}

	var v Value
	switch _type {
	case STRING:
		v = rdb.readStringValue()
	}
	v.Expire = expire

	return v
}

func (rdb *RDB) readTimestamp(b byte) int64 {
	var result int64
	switch b {
	case 0XFC:
		b := rdb.readBytes(8)
		d := binary.LittleEndian.Uint64(b)
		result = int64(d)
		 
	case 0XFD:
		b := rdb.readBytes(4)
		d := binary.LittleEndian.Uint32(b)
		result = int64(d)
	}
	
	return result
}

func (rdb *RDB) readStringValue() Value {
	keyName := rdb.readEncodedValue()
	value := rdb.readEncodedValue()

	return Value{
		Type: STRING,
		Key: keyName.Value,
		Value: value.Value,
	}
}

func (rdb *RDB) readByte() byte {
	b, err := rdb.Reader.ReadByte()
	if err != nil {
		panic(err)
	}

	return b
}

func (rdb *RDB) peek() byte {
	b, err := rdb.Reader.Peek(1)
	if err != nil {
		panic(err)
	}

	return b[0]
}

func (rdb *RDB) readBytes(n int) []byte {
	buf := make([]byte, n)
	n, err := rdb.Reader.Read(buf)
	if err != nil {
		panic(err)
	}

	return buf
}


func (rdb *RDB) readLenght() (int, bool) {
	var b byte
	b = rdb.peek()

	var lenght int
	var isSpecial bool
	switch b >> 6 {
	case 0b00:
		b := rdb.readByte()
		lenght = int(b)
	case 0b01:
		bytes := rdb.readBytes(2)
		firstByte := bytes[0]
		firstByte = firstByte & 0b00111111
		bytes[0] = firstByte

		value := binary.BigEndian.Uint16(bytes)

		lenght = int(value)
	case 0b10:
		rdb.readByte()
		bytes := rdb.readBytes(4)

		value := binary.BigEndian.Uint32(bytes)

		lenght = int(value)
	case 0b11:
		//rdb.readByte()
		return 0, true
		
	}
	
	return lenght, isSpecial
}

func (rdb *RDB) readEncodedValue() EncodedValue {
	lenght, isSpecial := rdb.readLenght()
	var value string
	if isSpecial {
		b := rdb.readByte()
		switch b {
		case 0xC0:
			v := int8(rdb.readByte())
            return EncodedValue{Value: strconv.Itoa(int(v))}
		case 0xC1:
			bytes := rdb.readBytes(2)
            v := int16(binary.LittleEndian.Uint16(bytes))
			return EncodedValue{Value: strconv.Itoa(int(v))}
		case 0xC2:
			bytes := rdb.readBytes(4)
            v := int32(binary.LittleEndian.Uint32(bytes))
            return EncodedValue{Value: strconv.Itoa(int(v))}
			case 0xC3:
				fmt.Printf("f")
		}
	} else {
		bytes := rdb.readBytes(lenght)
		value = string(bytes)
	}

	return EncodedValue{
		Value: value,
	}
}


