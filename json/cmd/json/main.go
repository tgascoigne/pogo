package main

import (
	"os"

	"github.com/tgascoigne/pogo/json"
)

func main() {
	path := os.Args[1]
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	json.Parse(file)
}
