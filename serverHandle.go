package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type KvStorage struct {
	mapval map[string]string
}

func NewKvStorage(port string) *KvStorage {
	return &KvStorage{
		mapval: make(map[string]string),
	}
}

func (k *KvStorage) DataStorage(instruction string, Path string) error {
	inst := strings.Fields(instruction)
	if len(inst) != 3 && (strings.ToLower(inst[0]) != "del" && inst[0] != "DELETE") {
		return fmt.Errorf("Invalid param")
	}
	update, err := NewWal(Path)
	if err != nil {
		return err
	}
	err = update.WriteFile(instruction)
	if err != nil {
		fmt.Println("Error while writing to file", err)
	}
	k.WriteToMemory(Path)
	return nil

}
func (k *KvStorage) WriteToMemory(filepath string) error {
	_, err := NewWal(":8080")
	if err != nil {
		return fmt.Errorf("Error in opening port")
	}
	fileOpen, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("Error in read file %s", err)
	}
	defer fileOpen.Close()
	scanner := bufio.NewScanner(fileOpen)
	for scanner.Scan() {
		entry := strings.Fields(scanner.Text())
		switch strings.ToLower(entry[0]) {
		case "set":
			k.mapval[entry[1]] = entry[2]
		case "del", "delete":
			val, ok := k.mapval[entry[1]]
			if !ok {
				return fmt.Errorf("Entry not found %s", val)
			} else {
				delete(k.mapval, entry[1])
			}
		}
	}
	return nil
}
