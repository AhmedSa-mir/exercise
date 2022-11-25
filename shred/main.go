package main

import (
	"fmt"
	"os"

	"example/pkg/shred"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage:", os.Args[0], "<file_path>")
		return
	}
	var path string = os.Args[1]
	err := shred.Shred(path)
	if err != nil {
		fmt.Println(err)
	}
}
