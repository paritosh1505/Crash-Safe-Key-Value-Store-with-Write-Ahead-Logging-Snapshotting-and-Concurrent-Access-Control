package main

import (
	"time"
)

func main() {
	//fileName := ManageWALFile()

	fileCurr := ReplayAllWAL("WAL_LOG")
	if fileCurr == "" {
		fileCurr = CreateFile(1)
	}
	rec1 := WALRecord{
		operation: "SET",
		key:       "key1",
		val:       "hello",
		timestamp: uint64(time.Now().UnixNano()),
	}
	rec2 := WALRecord{
		operation: "SET",
		key:       "key2",
		val:       "Bye",
		timestamp: uint64(time.Now().UnixNano()),
	}
	rec3 := WALRecord{
		operation: "DEL",
		key:       "key1",
		timestamp: uint64(time.Now().UnixNano()),
	}
	rec4 := WALRecord{
		operation: "SET",
		key:       "key3",
		val:       "this is 3rd key",
		timestamp: uint64(time.Now().UnixNano()),
	}
	rec5 := WALRecord{
		operation: "SET",
		key:       "key4",
		val:       "This is 4th key",
		timestamp: uint64(time.Now().UnixNano()),
	}
	rec6 := WALRecord{
		operation: "SET",
		key:       "key5",
		val:       "This is 5th key",
		timestamp: uint64(time.Now().UnixNano()),
	}
	rec7 := WALRecord{
		operation: "SET",
		key:       "key2",
		val:       "This is latest 2nd key",
		timestamp: uint64(time.Now().UnixNano()),
	}
	rec8 := WALRecord{
		operation: "SET",
		key:       "key3",
		val:       "This is latest 3rd key",
		timestamp: uint64(time.Now().UnixNano()),
	}

	WriteToWAL(fileCurr, rec1)
	WriteToWAL(fileCurr, rec2)
	WriteToWAL(fileCurr, rec3)
	WriteToWAL(fileCurr, rec4)
	WriteToWAL(fileCurr, rec5)
	WriteToWAL(fileCurr, rec6)
	WriteToWAL(fileCurr, rec7)
	WriteToWAL(fileCurr, rec8)

}
