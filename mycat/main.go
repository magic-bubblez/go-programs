package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <filename>")
	}

	for _, arg := range os.Args[1:] {
		file, err := os.Open(arg) //returns two val - []byte slice and error
		if err != nil {
			fmt.Println(os.Stderr, "mycat: %v\n", err)
			continue
		}
		_, err = io.Copy(os.Stdout, file)
		if err != nil {
			fmt.Println(os.Stderr, "mycat: error copying: %v\n", err)
		}
		file.Close()

	}
}
