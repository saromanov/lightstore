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
	_, err := w.file.Write(key)
	if err != nil {
		return fmt.Errorf("unable to write to file: %v", err)
	}
	return nil
}
