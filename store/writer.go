package store

import (
	"fmt"
	"os"
)

// Writer provides writing of data(operation log) to file
type Writer struct {
	file *os.File
}

// AddSetCommand defines append data to log file with SET command
func (w *Writer) AddSetCommand(key, value []byte) error {
	if err := w.write([]byte("set")); err != nil {
		return err
	}
	if err := w.write(key); err != nil {
		return err
	}
	if err := w.write(value); err != nil {
		return err
	}
	return w.write([]byte(";"))
}

func (w *Writer) write(data []byte) error {
	_, err := w.file.Write(data)
	if err != nil {
		return fmt.Errorf("unable to write to file: %v", err)
	}
	return nil
}
