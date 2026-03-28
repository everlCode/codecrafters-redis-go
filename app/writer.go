package main

import (
	"io"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type Writer struct {
	writer io.Writer
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{writer: writer}
}

func (w Writer) Write(v resp.Value) error {
	result := v.Marshal()
	_, err := w.writer.Write(result)

	return err
}
