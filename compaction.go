package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

type SnapShotFile struct {
	file *os.File
}
type SnapShotHeader struct {
	MagicNumber uint32
	Version     uint32
	EntryCount  uint32
}
type SnapShotEntry struct {
	Keylen uint32
	Vallen uint32
	Key    []byte
	Val    []byte
}

type SnapShotFooter struct {
	Crcval uint32
}

func CreateSnapShot() (*SnapShotFile, error) {
	file, err := os.OpenFile("snapshot.tmp", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal("Error while creating snapshot file ", err)
	}
	return &SnapShotFile{file}, nil
}
func DeleteSnapBinary() {
	fmt.Println("Deleting the Binary since we have encountered the error")
	err := os.Remove("snapshot.tmp")
	if err != nil {
		log.Fatal("Error while deleting the file ", err)
	}
}
func (h *SnapShotHeader) WriteSnapShotHeader(file *os.File) {
	if err := binary.Write(file, binary.BigEndian, h.MagicNumber); err != nil {
		DeleteSnapBinary()
		log.Fatal("Error while writing the file", err)

	}
	if err := binary.Write(file, binary.BigEndian, h.Version); err != nil {
		DeleteSnapBinary()
		log.Fatal("Error while writing version number ", err)
	}
	if err := binary.Write(file, binary.BigEndian, h.EntryCount); err != nil {
		DeleteSnapBinary()
		log.Fatal("Error while writing entry count ", err)
	}

}

func (h *SnapShotEntry) WriteSnapShotEntry() {
	//fmt.Println("****************", len(mapval))

}

func (h *SnapShotFooter) WriteSnapShotFooter(file *os.File) {
	if err := binary.Write(file, binary.BigEndian, h.Crcval); err != nil {
		DeleteSnapBinary()
		log.Fatal("Error while writng crc value ", err)
	}
}

func WriteToSnap() {
	headerval := SnapShotHeader{
		MagicNumber: magicNumber,
		Version:     uint32(versionNum),
		EntryCount:  uint32(len(mapval)),
	}
	entry := SnapShotEntry{}
	file, err := CreateSnapShot()
	if err != nil {
		log.Fatal("Error while creating in file ", err)
	}
	headerval.WriteSnapShotHeader(file.file)
	entry.WriteSnapShotEntry()
}
