package main

import (
	"bytes"
	"os"
	"sync"
)

const dirPath = "WAL_LOG"
const threshold = 109
const compaction_limit = 250
const magicNumber = 0x534E4150

var versionNum uint32 = 1

// var cummulative_size_compaction = 0

// var cummulative_buff_size int64 = 0
//var buff = new(bytes.Buffer)

//var sealedFile []string
//var sealedFileSet = make(map[string]struct{})

//var snap_path = "snapshot.tmp"
//var manifest = "manifest.json"

type CentralStorage struct {
	sealedFile                  []string
	sealedFileSet               map[string]struct{}
	cummulative_size_compaction int
	cummulative_buff_size       int64
	snap_path                   string
	buff                        *bytes.Buffer
	manifest                    string
	mapval                      map[string]string
	file                        *os.File
	mu                          sync.Mutex
}

var centralStorage = CentralStorage{
	sealedFile:                  make([]string, 0),
	sealedFileSet:               make(map[string]struct{}),
	cummulative_size_compaction: 0,
	cummulative_buff_size:       0,
	snap_path:                   "snapshot.tmp",
	manifest:                    "manifest.json",
	buff:                        new(bytes.Buffer),
	file:                        nil,
	mapval:                      make(map[string]string),
}
