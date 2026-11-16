package storage

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.OpenFile(
		"database.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666,
	)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	_, err = file.Write([]byte("SET a=10\n"))
	if err != nil {
		panic(err)
	}

	fmt.Println("Reccord append succesfully")

}
