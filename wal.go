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

func (c *CentralStorage) OpenWAL(path string) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	c.file = file
	c.currentOperationFile = file.Name()
	return nil
}
func WriteBinary(buffData *bytes.Buffer, data any) error {

	err := binary.Write(buffData, binary.BigEndian, data)
	if err != nil {
		return err
	}
	return nil
}

func (c *CentralStorage) EncodeWALRecord(r WALRecord) ([]byte, error, int64) {

	c.buff.Reset()
	var op byte

	if r.operation == "SET" {
		op = SET
	} else {
		op = DEL
	}

	err := WriteBinary(c.buff, op)
	if err != nil {
		return nil, err, 0
	}
	err = WriteBinary(c.buff, uint64(r.timestamp))
	if err != nil {
		return nil, err, 0
	}
	err = WriteBinary(c.buff, uint32(len(r.key)))
	if err != nil {
		return nil, err, 0
	}
	if op == 2 {
		r.val = ""
	}
	err = WriteBinary(c.buff, uint32(len(r.val)))
	if err != nil {
		return nil, err, 0
	}
	_, err = c.buff.Write([]byte(r.key))
	if err != nil {
		return nil, err, 0
	}
	_, err = c.buff.Write([]byte(r.val))
	if err != nil {
		return nil, err, 0
	}
	crcCalculate := crc32.ChecksumIEEE(c.buff.Bytes())
	err = WriteBinary(c.buff, crcCalculate)

	if err != nil {
		return nil, err, 0
	}
	c.cummulative_buff_size += int64(c.buff.Len())

	return c.buff.Bytes(), nil, c.cummulative_buff_size
}
func (c *CentralStorage) DecodeWALRecord(filereader *os.File) (WALRecord, error) {

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

func (c *CentralStorage) StoreSealedFile() []string {

	dirpath, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal("Error while reading directory in storeSealedFile")
	}

	for _, filename := range dirpath {
		if _, exist := c.sealedFileSet[filename.Name()]; !exist {
			c.sealedEntry = append(c.sealedEntry, filename.Name())
			c.sealedFileSet[filename.Name()] = struct{}{}
		}

	}
	defer c.file.Close()
	return c.sealedEntry
}

func (c *CentralStorage) WriteToWAL(currFile string, record WALRecord) error {
	c.mu.Lock()
	encodedrec, err, buffsize := c.EncodeWALRecord(record)
	defer c.mu.Unlock()
	info, err := os.Stat(currFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Scanning manifest json file now")
			currFile = c.ScanFileUsingManifest()
			info, _ = os.Stat(currFile)
			fmt.Println("curr file", currFile, "and size ", info.Size())

		} else {
			log.Fatal("Error in file operation-->", err)
		}

	}

	currFileSize := info.Size()
	currFileSize += int64(buffsize)
	fmt.Println("curr file", currFile, "and size ", currFileSize, " and sealed file ", c.sealedFile)

	c.cummulative_size_compaction += int(buffsize)
	if err != nil || buffsize > threshold || currFileSize > threshold {
		c.sealedFile = c.StoreSealedFile()
		c.currentOperationFile = c.ManageWALFile()
	}

	c.file.Write(encodedrec)
	c.file.Sync()
	if c.cummulative_size_compaction < compaction_limit {
		c.WriteToMap(record)
	} else {
		//TODO: implement background compaction
		c.TriggerCompaction()
		//c.cummulative_size_compaction = 0
	}

	return err
}
