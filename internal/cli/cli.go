package cli

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	flag "github.com/spf13/pflag"
)

type Choice int

const Version = "1.1.0"

const (
	Yes Choice = iota
	No
	Quit
	Dry
)

func (c Choice) String() string {
	switch c {
	case Yes:
		return "Yes"
	case No:
		return "No"
	case Quit:
		return "Quit"
	case Dry:
		return "Dry"
	default:
		return "Unknown" // Just here for future proofing if needed.
	}
}

type Arguments struct {
	Choice Choice
}

func Parse() ([]bool, Arguments) {
	var yes, no, quit, dry, version bool
	flag.BoolVarP(&yes, "yes", "y", false, "Launch in current terminal window")
	flag.BoolVarP(&no, "no", "n", false, "Don't launch in current terminal window")
	flag.BoolVarP(&quit, "quit", "q", false, "Quit before launching the application")
	flag.BoolVarP(&dry, "dry", "d", false, "Don't execute the selection, just print it")
	flag.BoolVarP(&version, "version", "v", false, "Print the current version")
	flag.Parse()

	selected := []bool{yes, no, dry, quit}
	if CountTrue(selected) > 1 {
		log.Fatal("You can only choose one from: [yes, no, quit]")
	}

	var choice Choice
	switch {
	case yes:
		choice = Yes
	case no:
		choice = No
	case quit:
		choice = Quit
	case dry:
		choice = Dry
	case version:
		fmt.Println("Tofi-menu version:", Version)
		os.Exit(0)
	}
	return selected, Arguments{Choice: choice}
}

func CountTrue(flags []bool) int {
	count := 0
	for _, b := range flags {
		if b {
			count++
		}
	}
	return count
}
