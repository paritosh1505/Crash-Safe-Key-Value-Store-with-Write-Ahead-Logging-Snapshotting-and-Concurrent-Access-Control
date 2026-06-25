package main

import (
	"time"
)

func main() {
	//var ct CentralStorage
	fileCurr := centralStorage.ReplayAllWAL("WAL_LOG")
	if fileCurr == "" {
		centralStorage.currentOperationFile = centralStorage.CreateWALSegment(1)
	}
	_ = centralStorage.OpenWAL(centralStorage.currentOperationFile)
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
	rec9 := WALRecord{
		operation: "SET",
		key:       "key6",
		val:       "Now i am adding 6th key",
		timestamp: uint64(time.Now().UnixNano()),
	}
	rec10 := WALRecord{
		operation: "SET",
		key:       "key8",
		val:       "hello eigth key",
		timestamp: uint64(time.Now().UnixNano()),
	}

	centralStorage.WriteToWAL(centralStorage.currentOperationFile, rec1)

	centralStorage.WriteToWAL(centralStorage.currentOperationFile, rec2)
	centralStorage.WriteToWAL(centralStorage.currentOperationFile, rec3)
	centralStorage.WriteToWAL(centralStorage.currentOperationFile, rec4)
	centralStorage.WriteToWAL(centralStorage.currentOperationFile, rec5)
	centralStorage.WriteToWAL(centralStorage.currentOperationFile, rec6)
	centralStorage.WriteToWAL(centralStorage.currentOperationFile, rec7)
	centralStorage.WriteToWAL(centralStorage.currentOperationFile, rec8)
	centralStorage.WriteToWAL(centralStorage.currentOperationFile, rec9)
	centralStorage.WriteToWAL(centralStorage.currentOperationFile, rec10)

}
