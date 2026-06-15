package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
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

func EncodeRecord(r WALRecord) ([]byte, error, int) {
	buff.Reset()
	var op byte

	if r.operation == "SET" {
		op = SET
	} else {
		op = DEL
	}

	err := WriteData(buff, op)
	if err != nil {
		return nil, err, 0
	}
	err = WriteData(buff, uint64(r.timestamp))
	if err != nil {
		return nil, err, 0
	}
	err = WriteData(buff, uint32(len(r.key)))
	if err != nil {
		return nil, err, 0
	}
	if op == 2 {
		r.val = ""
	}
	err = WriteData(buff, uint32(len(r.val)))
	if err != nil {
		return nil, err, 0
	}
	_, err = buff.Write([]byte(r.key))
	if err != nil {
		return nil, err, 0
	}
	_, err = buff.Write([]byte(r.val))
	if err != nil {
		return nil, err, 0
	}
	crcCalculate := crc32.ChecksumIEEE(buff.Bytes())
	err = WriteData(buff, crcCalculate)

	if err != nil {
		return nil, err, 0
	}

	cummulative_buff_size += buff.Len()

	return buff.Bytes(), nil, cummulative_buff_size
}
func DecodeRecord(filereader *os.File) (WALRecord, error) {

	var record WALRecord
	var operation byte
	var timestamp uint64
	var keylen uint32
	var vallen uint32
	var crc32 uint32
	err := binary.Read(filereader, binary.BigEndian, &operation)
	if err == io.EOF {
		return WALRecord{}, err
	}
	if err != nil {
		return WALRecord{}, err
	}
	switch operation {
	case 1:
		record.operation = "SET"
	case 2:
		record.operation = "DEL"
	default:
		return WALRecord{}, fmt.Errorf("Invalid operation %d", operation)
	}
	err = binary.Read(filereader, binary.BigEndian, &timestamp)
	if err != nil {
		return WALRecord{}, err
	}
	err = binary.Read(filereader, binary.BigEndian, &keylen)
	if err != nil {
		return WALRecord{}, err
	}
	err = binary.Read(filereader, binary.BigEndian, &vallen)
	if err != nil {
		return WALRecord{}, err
	}

	keylist := make([]byte, keylen)
	_, err = io.ReadFull(filereader, keylist)
	record.key = string(keylist)
	if err != nil {
		return WALRecord{}, err
	}
	vallist := make([]byte, vallen)
	_, err = io.ReadFull(filereader, vallist)
	record.val = string(vallist)
	if err != nil {
		return WALRecord{}, err
	}
	binary.Read(filereader, binary.BigEndian, &crc32)
	return record, nil

}

func StoreSealedFile(buffsize int, threshold int, sealedEntry []string) []string {
	dirpath, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal("Error while reading directory in storeSealedFile")
	}

	for _, filename := range dirpath {
		if _, exist := sealedFileSet[filename.Name()]; !exist {
			sealedEntry = append(sealedEntry, filename.Name())
			sealedFileSet[filename.Name()] = struct{}{}
		}

	}
	return sealedEntry
}
func WriteToWAL(currFile string, record WALRecord) error {
	var filepath string
	encodedrec, err, buffsize := EncodeRecord(record)

	info, err := os.Stat(currFile)
	if err != nil {
		log.Fatal("Getting error in stat function-->", err)
	}
	currFileSize := info.Size()
	currFileSize += int64(buffsize)
	//fmt.Println("**************", sealedFile)
	cummulative_size_compaction += int(buffsize)
	if err != nil || buffsize > threshold || currFileSize > threshold {
		sealedFile = StoreSealedFile(buffsize, threshold, sealedFile)
		filepath = ManageWALFile()
	} else {
		filepath = currFile
	}
	walBuff, err := NewWal(filepath)
	if err != nil {
		return err
	}
	defer walBuff.file.Close()

	walBuff.file.Write(encodedrec)
	walBuff.file.Sync()
	//record, err = DecodeRecord()
	if cummulative_size_compaction < compaction_limit {
		WriteToMap(record)
	} else {
		//sealedFile = append(sealedFile, currFile)
		WriteToSnap()
		cummulative_size_compaction = 0
	}

	return err
}
