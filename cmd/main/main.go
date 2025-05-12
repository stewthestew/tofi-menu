package main

import (
	"fmt"
	"log"
	"os"
	"tofi/internal/backend"
)

func main() {
	list, err := backend.List(os.Getenv("PATH"))
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range list {
		fmt.Println(v)
	}
}
