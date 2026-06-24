package main

import (
	"time"
)

func main1() {
	fileCurr := ReplayAllWAL("WAL_LOG")
	if fileCurr == "" {
		fileCurr = CreateWALSegment(1)
	}
	c1, _ := OpenWAL(fileCurr)
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

	c1.WriteToWAL(fileCurr, rec1)
	c1.WriteToWAL(fileCurr, rec2)
	c1.WriteToWAL(fileCurr, rec3)
	c1.WriteToWAL(fileCurr, rec4)
	c1.WriteToWAL(fileCurr, rec5)
	c1.WriteToWAL(fileCurr, rec6)
	c1.WriteToWAL(fileCurr, rec7)
	c1.WriteToWAL(fileCurr, rec8)
	c1.WriteToWAL(fileCurr, rec9)
	c1.WriteToWAL(fileCurr, rec10)

}
