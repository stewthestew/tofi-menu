package tofiLib

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"tofi/internal/backend"
	"tofi/internal/cli"

	"github.com/charmbracelet/log"
)

const DefaultLaunchApp = "fzf"

type TofiOptions struct {
	SelectedArgs []bool
	Args         cli.Arguments
	List         []string
	Silent       bool
}

type EnvVars struct {
	Debug     bool
	LaunchApp string
}

func TofiOptionsInit() TofiOptions {
	selectedArgs, args := cli.Parse()
	silent := os.Getenv("SILENT") != ""

	list, err := backend.List(os.Getenv("TOFI_APPS"), silent)
	if err != nil {
		log.Fatal(err)
	}

	return TofiOptions{
		SelectedArgs: selectedArgs,
		Args:         args,
		List:         list,
		Silent:       silent,
	}
}

func EnvVarsInit(silent bool) EnvVars {
	launchApp := os.Getenv("LAUNCH_APP")
	if launchApp == "" {
		log.Warn("Could not find fzf in FZF_PATH. If the environment variable FZF_PATH is not set to the fzf binary, please set it to the fzf binary.")
		log.Warn("Trying to run it anyway...")
		launchApp = DefaultLaunchApp
	}

	return EnvVars{
		Debug:     os.Getenv("DEBUG") != "",
		LaunchApp: launchApp,
	}
}

func Execute(o TofiOptions, e EnvVars, executedCmd *exec.Cmd, selectedCmd string) error {
	var choice string
	if cli.CountTrue(o.SelectedArgs) == 0 {
		reader := strings.NewReader("yes\nno")
		if !o.Silent {
			user := exec.Command(e.LaunchApp)
			user.Stdin = reader
			tempChoice, err := user.Output()
			if err != nil {
				return errors.New(fmt.Sprintf("Error reading launchApp output %v", err))
			}
			choice = strings.TrimSpace(string(tempChoice))
		}

		choice = strings.TrimSpace(choice)
	} else {
		choice = o.Args.Choice.String()
	}

	// if cli.CountTrue(selectedArgs) == 0 {
	// 	user := exec.Command(fzf)
	// 	user.Stdin = strings.NewReader("yes\nno")
	// 	choiceb, err := user.Output()
	// 	if err != nil {
	// 		log.Error(err)
	// 	}
	//
	// 	choice = strings.TrimSpace(string(choiceb))
	// } else {
	// 	choice = args.Choice.String()
	// }

	log.Debug("User selected:", choice)

	switch strings.ToLower(choice) {

	case "n", "no":
		executedCmd.SysProcAttr = &syscall.SysProcAttr{
			Setsid: true,
		}
		err := executedCmd.Start()
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to Start process %v", err))
		}

	case "d", "dry":
		fmt.Print(selectedCmd)
	case "y", "yes":
		executedCmd.Stdout = os.Stdout
		executedCmd.Stdin = os.Stdin
		executedCmd.Stderr = os.Stderr

		err := executedCmd.Run()
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to Run process %v", err))
		}

	case "q", "quit", "":
		log.Info("Bye!")

	default:
		return errors.New(fmt.Sprintf("Invalid argument %v", choice))
	}

	return nil
}

func Tofi(o TofiOptions, e EnvVars) {
	if o.Silent {
		log.SetOutput(io.Discard)
	}

	if e.Debug {
		log.SetLevel(log.DebugLevel)
	}

	a := strings.Join(o.List, "\n")

	cmd := exec.Command(e.LaunchApp)
	cmd.Stdin = strings.NewReader(a)
	selected, err := cmd.Output()
	// The reason for this is: I assuming am if the program exits with an error with 0 or 1 it errors out. So I am only handling the error if the debug env var is set.
	// This pattern will be used throughout the code.
	// ^I turned out to be kind of right, apparently fzf when exiting without selecting something it exists with exit code "130"
	if err != nil {
		log.Debug("Fzf errored out", err)
	}

	selectedCmd := strings.TrimSpace(string(selected))

	if selectedCmd == "" {
		log.Info("No command selected. Exiting...")
		return
	}

	log.Debug("selected", selectedCmd)

	fields := strings.Fields(selectedCmd)
	executedCmd := exec.Command(fields[0], fields[1:]...)
	// .RUN() Is blocking which is good for tui apps, but for gui's? No.

	err = Execute(o, e, executedCmd, selectedCmd)
	if err != nil {
		log.Fatal(err)
	}
}
