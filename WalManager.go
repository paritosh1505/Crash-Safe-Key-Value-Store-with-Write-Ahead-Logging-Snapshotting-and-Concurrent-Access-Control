package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func (c *CentralStorage) ManageWALFile() string {

	dir, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatalf("Directory %s not found: %v", dirPath, err)
	}

	maxIndex := 0
	var latestFileSize int64
	var latestFileName string

	// 1. Find only the highest index and its size
	for _, entry := range dir {
		name := entry.Name()
		if strings.HasPrefix(name, "wal-") && strings.HasSuffix(name, ".log") {
			trimmed := strings.TrimSuffix(strings.TrimPrefix(name, "wal-"), ".log")
			idx, err := strconv.Atoi(trimmed)
			if err == nil && idx > maxIndex {
				maxIndex = idx
				latestFileName = name
				info, _ := entry.Info()
				latestFileSize = info.Size()
			}
		}
	}
	latestFileSize += int64(c.cummulative_buff_size)
	var newfile string
	if maxIndex == 0 {
		newfile = c.CreateWALSegment(1)
	} else if latestFileSize < int64(threshold) {
		newfile = dirPath + "/" + latestFileName
	} else {
		maxIndex = maxIndex + 1
		newfile = c.CreateWALSegment(maxIndex)

	}
	c.buff.Reset()
	c.cummulative_buff_size = 0
	return newfile
}
func (c *CentralStorage) CreateWALSegment(index int) string {
	path := fmt.Sprintf("WAL_LOG/wal-%06d.log", index)
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			c.file, err = os.Create(path)
			if err != nil {
				log.Fatal("Error while creating teh file")
			}
		} else {
			log.Fatal("Unexpected error while creating the file")
		}
	}
	fmt.Println(path)
	return path
}

func (c *CentralStorage) ReplayAllWAL(dirName string) string {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, err := os.Stat(c.snap_path)
	if err == nil {
		err = c.LoadSnapDataToMemory()
		if err != nil {
			fmt.Println("Error in snapshot file-->", err)
		}
	}
	index := c.FetchManifestIndex()
	var currFileName string
	dirpath, err := os.ReadDir(dirName)
	if err != nil {
		log.Fatal("Directory Not found")
	}
	//index:=FetchManifestIndex()
	for i, filename := range dirpath {
		if i > index {
			currFileName = dirName + "/" + filename.Name()
			filereader, err := os.Open(currFileName)
			if err != nil {
				log.Fatal("Error while reading the file")
			}
			for {
				record, err := c.DecodeWALRecord(filereader)
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Printf("Error decoding WAL record in %s: %v", currFileName, err)
					break
				}
				c.WriteToMap(record)
			}

		}

		if len(currFileName) == 0 {
			return ""
		} else {
			stat, err := os.Stat(currFileName)
			fmt.Println("Curr file name and size", currFileName, stat.Size())
			if err != nil {
				log.Fatal("Error while doing stat reading of file")
			}
			if stat.Size() < threshold {
				return currFileName
			}
		}
	}

	return ""
}
