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
var snapshotMap map[string]string

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
	manifestJson                ManifestJson
	currentOperationFile        string
	isCompactionRunning         bool
	sealedEntry                 []string
	mu                          sync.Mutex
}
type ManifestJson struct {
	Manifest_version   int    `json:"version"`
	Last_compact_index int    `json:"last_index"`
	Snapshot_file      string `json:"file"`
	Snapshot_checksum  uint32 `json:"checksum"`
	Created_at         int64  `json:"CreatedAt"`
	Entry_count        int    `json:"entryCount"`
}

var centralStorage = CentralStorage{
	sealedFile:                  make([]string, 0),
	sealedFileSet:               make(map[string]struct{}),
	cummulative_size_compaction: 0,
	cummulative_buff_size:       0,
	currentOperationFile:        "",
	snap_path:                   "snapshot.tmp",
	manifest:                    "manifest.json",
	buff:                        new(bytes.Buffer),
	file:                        nil,
	manifestJson:                ManifestJson{},
	mapval:                      make(map[string]string),
}

type SnapShotFile struct {
	file *os.File
}
type ManifestFile struct {
	file *os.File
}
type SnapShotHeader struct {
	MagicNumber uint32
	Version     uint32
	EntryCount  uint32
}
type SnapShotEntry struct {
	Keylen int32
	Vallen int32
	Key    string
	Val    string
}

type SnapShotFooter struct {
	Crcval uint32
}
