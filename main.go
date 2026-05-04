package main

import "fmt"

func main() {
	store := NewKvStorage("8080")
	_ = store.DataStorage("SET key1 hello", "WAL.log")
	_ = store.DataStorage("Del key1", "WAL.log")
	_ = store.DataStorage("SET key2 hello", "WAL.log")

	for k, v := range store.mapval {
		fmt.Println(k, "--->", v)
	}

}
