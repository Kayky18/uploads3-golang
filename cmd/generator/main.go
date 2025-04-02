package main

import (
	"fmt"
	"os"
)

func main() {
	i := 0
	for {
		f, err := os.Create(fmt.Sprintf("./tmp/file%d.txt", i))
		if err != nil {
			fmt.Println("Error creating file:", err)
			break
		}
		defer f.Close()
		_, err = f.WriteString("Hello, World!\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			break
		}
		i++
		if i >= 10 {
			break
		}
	}
}
