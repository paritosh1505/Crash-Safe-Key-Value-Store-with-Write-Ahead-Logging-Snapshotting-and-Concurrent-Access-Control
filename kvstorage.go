package main

var mapval = make(map[string]string)

func (c *CentralStorage) WriteToMap(record WALRecord) {

	switch record.operation {
	case "SET":
		c.mapval[record.key] = record.val
	case "DEL":
		delete(c.mapval, record.key)
	}
}
