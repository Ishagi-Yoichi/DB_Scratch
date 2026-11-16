package main

import (
	"fmt"
	"os"
)

func main() {
	data := []byte("SET a=10\n")

	fmt.Println("Original data:", data)
	fmt.Println("As text:", string(data))

	for i, b := range data {
		fmt.Printf("Byte %d: %d\n", i, b)
	}

	file, err := os.OpenFile(
		"database.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666,
	)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		panic(err)
	}
}
