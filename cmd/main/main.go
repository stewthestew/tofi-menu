package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"tofi/internal/backend"
)

func main() {
	list, err := backend.List(os.Getenv("PATH"))
	fzf := os.Getenv("FZF_PATH")

	if len(os.Getenv("FZF_PATH")) == 0 {
		fmt.Println("Could not find fzf in FZF_PATH. If the environment variable FZF_PATH is not set to the fzf binary, please set it to the fzf binary.")
		fmt.Println("Trying to run it anyway...")
		fzf = "fzf"
	}

	if err != nil {
		log.Fatal(err)
	}

	a := strings.Join(list, "\n")

	cmd := exec.Command(fzf)
	cmd.Stdout = os.Stdout
	cmd.Stdin = strings.NewReader(a)
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		if len(os.Getenv("DEBUG")) > 0 {
			log.Fatal(err)
		}
	}

	selected, err := cmd.Output()
	if err != nil {
		if len(os.Getenv("DEBUG")) > 0 {
			log.Fatal(err)
		}
	}

	selectedCmd := strings.TrimSpace(string(selected))

	if selectedCmd == "" {
		return
	}

	executedCmd := exec.Command(selectedCmd)
	executedCmd.Stdout = os.Stdout
	executedCmd.Stderr = os.Stderr
	err = executedCmd.Run()
	if err != nil {
		if len(os.Getenv("DEBUG")) > 0 {
			log.Fatal(err)
		}
	}
}
