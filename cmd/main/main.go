package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"tofi/internal/backend"
)

func main() {
	list, err := backend.List(os.Getenv("PATH"))
	if err != nil {
		log.Fatal(err)
	}

	a := strings.Join(list, "\n")

	cmd := exec.Command("/bin/sh", "-c", a)
	cmd.Stdout = os.Stdout
	cmd.Stdin = strings.NewReader(a)
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
