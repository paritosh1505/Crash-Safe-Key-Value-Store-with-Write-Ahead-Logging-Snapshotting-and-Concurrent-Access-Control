package main

var mapval = make(map[string]string)

func WriteToMap(record WALRecord) {
	switch record.operation {
	case "SET":
		centralStorage.mapval[record.key] = record.val
	case "DEL":
		delete(centralStorage.mapval, record.key)
	}
}
