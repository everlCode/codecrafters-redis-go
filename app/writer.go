package main

import "io"

type Writer struct {
	writer io.Writer
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{writer: writer}
}

func (w Writer) Write(v Value) error {
	result := v.Marshal()
	_, err := w.writer.Write(result)

	return err
}
