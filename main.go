package main

import (
	"time"
)

func main() {
	rec1 := WALRecord{
		operation: "SET",
		key:       "key1",
		val:       "hello",
		timestamp: uint64(time.Now().Unix()),
	}
	rec2 := WALRecord{
		operation: "SET",
		key:       "key2",
		val:       "hello",
		timestamp: uint64(time.Now().Unix()),
	}
	WriteToWAL("WAL.log", rec1)
	WriteToWAL("WAL.log", rec2)

}
