package main

import (
	"bytes"
)

const dirPath = "WAL_LOG"
const threshold = 109
const compaction_limit = 250
const magicNumber = 0x534E4150

var cummulative_size_compaction = 0
var versionNum uint32 = 1
var cummulative_buff_size = 0
var buff = new(bytes.Buffer)

var sealedFile []string
var sealedFileSet = make(map[string]struct{})
var snap_path = "snapshot.tmp"
var manifest = "manifest.json"
