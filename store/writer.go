package store

import (
	"fmt"
	"os"
)

const defaultWritePath = "ligthstore.db"

// Writer provides writing of data(operation log) to file
type Writer struct {
	file *os.File
}

// newWriter provides initialization of the Writer
func newWriter(path string) (*Writer, error) {
	if path == "" {
		path = defaultWritePath
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	return &Writer{
		file: f,
	}, nil
}

// AddSetCommand defines append data to log file with SET command
func (w *Writer) AddSetCommand(key, value []byte) error {
	if err := w.write([]byte("set;\n")); err != nil {
		return err
	}
	if err := w.write(key); err != nil {
		return err
	}
	if err := w.write([]byte("\n")); err != nil {
		return err
	}
	if err := w.write(value); err != nil {
		return err
	}
	return w.write([]byte("\nend;\n"))
}

func (w *Writer) write(data []byte) error {
	_, err := w.file.Write(data)
	if err != nil {
		return fmt.Errorf("unable to write to file: %v", err)
	}
	return nil
}
