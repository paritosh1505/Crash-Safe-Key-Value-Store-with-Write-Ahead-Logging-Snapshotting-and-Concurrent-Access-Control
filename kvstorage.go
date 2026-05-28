package main

import "fmt"

var mapval = make(map[string]string)

func WriteToMap(record WALRecord) {
	switch record.operation {
	case "SET":
		mapval[record.key] = record.val
	case "DEL":
		delete(mapval, record.key)
	}
	fmt.Printf("%p\n", mapval)
}
