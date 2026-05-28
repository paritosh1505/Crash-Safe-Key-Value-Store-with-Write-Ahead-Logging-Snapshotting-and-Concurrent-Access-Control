package main

import (
	"fmt"
	"hash/crc32"
	"os"
	"strconv"
	"strings"
)

type Wal struct {
	file *os.File
}
type WALRecord struct {
	operation string
	key       string
	val       string
	timestamp uint64
	crcval    uint64
}

func NewWal(path string) (*Wal, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return &Wal{file}, nil
}

func EncodeRecord(r WALRecord) []byte {
	payload := fmt.Sprintf("%s|%s|%s|%d",
		r.operation,
		r.key,
		r.val,
		r.timestamp,
	)
	r.crcval = uint64(crc32.ChecksumIEEE([]byte(payload)))
	return []byte(fmt.Sprintf("%s|%d", payload, r.crcval))

}
func DecodeRecord(newrec []byte) WALRecord {
	trimNewLine := strings.TrimSpace(string(newrec))
	recArray := strings.Split(string(trimNewLine), "|")
	var record WALRecord
	record.operation = recArray[0]
	record.key = recArray[1]
	record.val = recArray[2]
	timestamp, err := strconv.ParseUint(recArray[3], 10, 64)
	if err != nil {
		fmt.Println("err while converting the timestamp")
		return WALRecord{}
	}
	record.timestamp = uint64(timestamp)
	payloadNew := fmt.Sprintf("%s|%s|%s|%d",
		record.operation,
		record.key,
		record.val,
		record.timestamp)
	crc, err := strconv.ParseUint(recArray[4], 10, 64)
	if err != nil {
		fmt.Println("Crc error")
		return WALRecord{}
	}
	record.crcval = uint64(crc32.ChecksumIEEE([]byte(payloadNew)))
	if record.crcval != crc {
		fmt.Println("CRC mismatch error", err)
		return WALRecord{}
	}
	return record

}

func WriteToWAL(walPath string, record WALRecord) error {
	wal, err := NewWal(walPath)
	if err != nil {
		return err
	}
	defer wal.file.Close()
	newrec := EncodeRecord(record)
	newrec = append(newrec, '\n')
	_, err = wal.file.Write(newrec)
	err = wal.file.Sync()
	if err != nil {
		fmt.Println("Sync fail")
		return err
	}
	recval := DecodeRecord(newrec)
	fmt.Println(recval)
	//WriteToMap(record)
	return err
}
