package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"tofi/internal/backend"

	"github.com/charmbracelet/log"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	list, err := backend.List(os.Getenv("TOFI_APPS"))
	if err != nil {
		log.Fatal(err)
	}

	fzf := os.Getenv("FZF_PATH")
	var DEBUG bool

	if len(os.Getenv("DEBUG")) > 0 {
		DEBUG = true
	}

	if len(os.Getenv("FZF_PATH")) == 0 {
		fmt.Println("Could not find fzf in FZF_PATH. If the environment variable FZF_PATH is not set to the fzf binary, please set it to the fzf binary.")
		fmt.Println("Trying to run it anyway...")
		fzf = "fzf"
	}

	a := strings.Join(list, "\n")

	cmd := exec.Command(fzf)
	cmd.Stdin = strings.NewReader(a)
	selected, err := cmd.Output()

	// The reason for this is: I assuming am if the program exits with an error with 0 or 1 it errors out. So I am only handling the error if the debug env var is set.
	// This pattern will be used throughout the code.
	// ^I turned out to be kind of right, apparently fzf when exiting without selecting something it exists with exit code "130"
	if err != nil {
		if DEBUG {
			log.Error("Fzf errored out", err)
		}
	}

	selectedCmd := strings.TrimSpace(string(selected))

	if selectedCmd == "" {
		log.Info("No command selected. Exiting...")
		return
	}

	if DEBUG {
		log.Debug("selected", selectedCmd)
	}

	executedCmd := exec.Command(selectedCmd)
	// .RUN() Is blocking which is good for tui apps, but for gui's? No.
	fmt.Print("Do you want to run this in the current window? [yes\\no\\Quit]\n: ")

	var choice string
	choice, err = reader.ReadString('\n')
	if err != nil {
		log.Error(err)
	}

	choice = strings.TrimSpace(choice)

	switch strings.ToLower(choice) {

	case "n", "no":
		executedCmd.Stdout = nil
		executedCmd.Stderr = nil
		executedCmd.Stdin = nil
		executedCmd.SysProcAttr = &syscall.SysProcAttr{
			Setsid: true,
		}
		err = executedCmd.Start()
	case "y", "yes":
		executedCmd.Stdout = os.Stdout
		executedCmd.Stdin = os.Stdin
		executedCmd.Stderr = os.Stderr

		err = executedCmd.Run()

	default:
		log.Info("Bye!")

	}

	if err != nil {
		log.Error(err)
	}

}
