package main

import (
	"fmt"
	"hash/crc32"
	"os"
	"time"
)

type WAL struct {
	file *os.File
}

func (w *WAL) Close() {
	w.file.Close()
}
func NewWal(path string) (*WAL, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &WAL{file: f}, nil
}

func (w *WAL) WriteFile(data string) error {

	entry := data + " " + time.Now().UTC().Format(time.RFC3339)
	crcData := crc32.ChecksumIEEE([]byte(entry))
	final := fmt.Sprintf("%s,%d",
		entry, crcData)
	_, err := w.file.Write([]byte(final + "\n"))
	if err != nil {
		return fmt.Errorf("Error while returning the file %s", err)
	}
	defer w.Close()
	return nil

}
