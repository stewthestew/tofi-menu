package tofiLib

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"tofi/internal/backend"
	"tofi/internal/cli"
	"tofi/internal/utils"

	"github.com/charmbracelet/log"
)

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
	silent := os.Getenv("TOFI_SILENT") != ""

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
	var launchApp string

	if len(os.Getenv("LAUNCH_APP")) == 0 {
		utils.Log(log.WarnLevel, silent, "Could not find fzf in FZF_PATH. If the environment variable FZF_PATH is not set to the fzf binary, please set it to the fzf binary.")
		utils.Log(log.WarnLevel, silent, "Trying to run it anyway...")
		launchApp = "fzf"
	}

	return EnvVars{
		Debug:     os.Getenv("DEBUG") != "",
		LaunchApp: launchApp,
	}
}

func Tofi(opts TofiOptions, vars EnvVars) {
	if vars.Debug {
		log.SetLevel(log.DebugLevel)
	}

	if len(os.Getenv("LAUNCH_APP")) == 0 {
		utils.Log(log.WarnLevel, opts.Silent, "Could not find fzf in FZF_PATH. If the environment variable FZF_PATH is not set to the fzf binary, please set it to the fzf binary.")
		utils.Log(log.WarnLevel, opts.Silent, "Trying to run it anyway...")
		vars.LaunchApp = "fzf"
	}

	a := strings.Join(opts.List, "\n")

	cmd := exec.Command(vars.LaunchApp)
	cmd.Stdin = strings.NewReader(a)
	selected, err := cmd.Output()

	// The reason for this is: I assuming am if the program exits with an error with 0 or 1 it errors out. So I am only handling the error if the debug env var is set.
	// This pattern will be used throughout the code.
	// ^I turned out to be kind of right, apparently fzf when exiting without selecting something it exists with exit code "130"
	if err != nil {
		utils.Log(log.DebugLevel, opts.Silent, "Fzf errored out", err)
	}

	selectedCmd := strings.TrimSpace(string(selected))

	if selectedCmd == "" {
		utils.Log(log.InfoLevel, opts.Silent, "No command selected. Exiting...")
		return
	}

	utils.Log(log.DebugLevel, opts.Silent, "selected", selectedCmd)

	fields := strings.Fields(selectedCmd)
	executedCmd := exec.Command(fields[0], fields[1:]...)
	// .RUN() Is blocking which is good for tui apps, but for gui's? No.

	var choice string
	if cli.CountTrue(opts.SelectedArgs) == 0 {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Do you want to run this in the current window? [yes\\no\\dry\\Quit]\n: ")

		choice, err = reader.ReadString('\n')
		if err != nil {
			utils.Log(log.ErrorLevel, opts.Silent, err)
		}

		choice = strings.TrimSpace(choice)
	} else {
		choice = opts.Args.Choice.String()
	}

	utils.Log(log.DebugLevel, opts.Silent, "User selected:", choice)

	switch strings.ToLower(choice) {

	case "n", "no":
		executedCmd.Stdout = nil
		executedCmd.Stderr = nil
		executedCmd.Stdin = nil
		executedCmd.SysProcAttr = &syscall.SysProcAttr{
			Setsid: true,
		}
		err = executedCmd.Start()

	case "d", "dry":
		fmt.Print(selectedCmd)
	case "y", "yes":
		executedCmd.Stdout = os.Stdout
		executedCmd.Stdin = os.Stdin
		executedCmd.Stderr = os.Stderr

		err = executedCmd.Run()
	case "q", "quit", "":
		utils.Log(log.InfoLevel, opts.Silent, "Bye!")

	default:
		log.Fatalf("Invalid argument")
		utils.Log(log.FatalLevel, opts.Silent, "Invalid argument")

	}

	if err != nil {
		utils.Log(log.ErrorLevel, opts.Silent, err)
	}

}
