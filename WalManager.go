package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func ManageWALFile() string {

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
	latestFileSize += int64(cummulative_buff_size)
	var newfile string
	if maxIndex == 0 {
		newfile = CreateFile(1)
	} else if latestFileSize < int64(threshold) {
		newfile = dirPath + "/" + latestFileName
	} else {
		maxIndex = maxIndex + 1
		newfile = CreateFile(maxIndex)
	}
	buff.Reset()
	cummulative_buff_size = 0
	return newfile
}
func CreateFile(index int) string {
	path := fmt.Sprintf("WAL_LOG/wal-%06d.log", index)
	fmt.Println(path)
	file, err := os.Create(path)
	if err != nil {
		log.Fatal("Error while creating teh file")
	}
	defer file.Close()
	return path
}

func ReplayAllWAL(dirName string) string {
	var currFileName string
	dirpath, err := os.ReadDir(dirName)
	if err != nil {
		log.Fatal("Directory Not found")
	}
	//index:=FetchManifestIndex()
	for _, filename := range dirpath {
		currFileName = dirName + "/" + filename.Name()
		filereader, err := os.Open(currFileName)
		if err != nil {
			log.Fatal("Error while reading the file")
		}
		for {
			record, err := DecodeRecord(filereader)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Error decoding WAL record in %s: %v", currFileName, err)
				break
			}
			WriteToMap(record)
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

	return ""
}
