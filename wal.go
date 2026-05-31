package main

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"io"
	"log"
	"os"
)

type Wal struct {
	file *os.File
}

const (
	SET byte = 0x01
	DEL byte = 0x02
)

type WALRecord struct {
	operation string
	key       string
	val       string
	timestamp uint64
	crcval    uint64
}

var op byte

func NewWal(path string) (*Wal, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return &Wal{file}, nil
}
func WriteData(buffData *bytes.Buffer, data any) error {

	err := binary.Write(buffData, binary.BigEndian, data)
	if err != nil {
		return err
	}
	return nil
}

func EncodeRecord(r WALRecord) ([]byte, error) {
	var buff = new(bytes.Buffer)
	if r.operation == "SET" {
		op = SET
	} else {
		op = DEL
	}

	err := WriteData(buff, op)
	if err != nil {
		return nil, err
	}
	err = WriteData(buff, uint64(r.timestamp))
	if err != nil {
		return nil, err
	}
	err = WriteData(buff, uint32(len(r.key)))
	if err != nil {
		return nil, err
	}
	if op == 2 {
		r.val = ""
	}
	err = WriteData(buff, uint32(len(r.val)))
	if err != nil {
		return nil, err
	}
	_, err = buff.Write([]byte(r.key))
	if err != nil {
		return nil, err
	}
	_, err = buff.Write([]byte(r.val))
	if err != nil {
		return nil, err
	}
	crcCalculate := crc32.ChecksumIEEE(buff.Bytes())
	err = WriteData(buff, crcCalculate)

	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}
func DecodeRecord() (WALRecord, error) {
	filereader, err := os.Open("WAL.log")
	if err != nil {
		log.Fatal("invalid file read", err)
	}
	var record WALRecord
	var operation byte
	var timestamp uint64
	var keylen uint32
	var vallen uint32
	var crc32 uint32
	binary.Read(filereader, binary.BigEndian, &operation)
	if operation == 1 {
		record.operation = "SET"
	} else {
		record.operation = "DEL"
	}
	err = binary.Read(filereader, binary.BigEndian, &timestamp)
	if err != nil {
		log.Fatal("error while reading timestamp==>", err)
	}
	err = binary.Read(filereader, binary.BigEndian, &keylen)
	if err != nil {
		log.Fatal("Error while reading binary key->", err)
	}
	err = binary.Read(filereader, binary.BigEndian, &vallen)
	if err != nil {
		log.Fatal("Error while reading binary val->", err)
	}

	keylist := make([]byte, keylen)
	_, err = io.ReadFull(filereader, keylist)
	record.key = string(keylist)
	if err != nil {
		log.Fatal("Error while reading key value->", keylist)
	}
	vallist := make([]byte, vallen)
	_, err = io.ReadFull(filereader, vallist)
	record.val = string(vallist)
	if err != nil {
		log.Fatal("Error while reading key value->", vallist)
	}
	binary.Read(filereader, binary.BigEndian, &crc32)
	return record, nil

}

func WriteToWAL(walPath string, record WALRecord) error {

	walBuff, err := NewWal(walPath)
	if err != nil {
		return err
	}
	defer walBuff.file.Close()
	encodedrec, err := EncodeRecord(record)
	if err != nil {
		log.Fatal("Hey there is error->", err)
	}
	walBuff.file.Write(encodedrec)
	walBuff.file.Sync()
	//record, err = DecodeRecord()
	WriteToMap(record)
	return err
}
