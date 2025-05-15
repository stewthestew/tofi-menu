![Go Workflow](https://github.com/stewthestew/tofi-menu/actions/workflows/go.yml/badge.svg)
# Tofi
Tofi or Tofi-menu is a simple app launcher for the terminal

## What is it? 
Tofi-menu is a basic wrapper around fzf, which allows you to launch applications from... Guess what? FZF!

> [!NOTE]
> Tofi is supposed to be used with Ghostty (Ghostty is a terminal emulator if you didn't know)

## Usage

Ghostty:
```bash
ghostty -e "<path to tofi>"
```

### Installation

```bash
git clone https://github.com/stewthestew/tofi-menu
cd tofi-menu
go build -o tofi cmd/main/main.go
mv ./tofi ~/go/bin/ # Or wherever you want to put it
```
### Important info
Tofi-menu uses 3 different environment variables:

1. FZF_PATH: The path to the fzf binary. (Optional, if not set then it will just execute fzf normally)
2. DEBUG: If this is set, then less important errors will be handled (Read code for more info) (Optional)
3. TOFI_APPS: This acts like the PATH variable, this is where tofi will get the executables from.

If you want it to list every app then do 
```bash
export TOFI_APPS=$PATH
```

You can also change TOFI_APPS to point at specific directories so you can launch custom scripts, etc.

## Diagram

I was too lazy to make my own diagram of how tofi works so this should be good enough
```text
main()
│
├── selectedArgs, args := cli.Parse()                         // Parse CLI flags: --yes, --no, --quit
│
├── list, err := backend.List($TOFI_APPS)                     // List all executable files from PATH
│
├── fzf := $FZF_PATH or fallback to "fzf"                     // Determine fzf binary path
│
├── cmd := exec.Command(fzf)                                  // Prepare to run fzf
│   └── cmd.Stdin = strings.NewReader(appsList)               // Feed app list to fzf
│
├── selected := cmd.Output()                                  // Get selection from fzf
│
├── selectedCmd := strings.TrimSpace(selected)                // Sanitize selected command
│   └── if selectedCmd == "" → log.Info + return
│
├── fields := strings.Fields(selectedCmd)                     // Split into command + args
├── executedCmd := exec.Command(fields[0], fields[1:]...)     // Prepare command execution
│
├── if cli.CountTrue(selectedArgs) == 0                       // If no CLI flags, prompt user
│   └── prompt "Do you want to run this in the current window?"
│
├── switch strings.ToLower(choice)                            // Handle user/flag choice
│   ├── "no"  → run in background (detached)
│   ├── "yes" → run blocking in same terminal
│   ├── "quit"/"" → exit
│   └── default → log.Fatalf
│
└── if err != nil → log.Error()


List(PATH string) → []string
│
├── Split PATH by ":"                                          // Multiple directories
├── For each path:
│   └── walkDir(path)
│       ├── Recursively walk files
│       ├── For each file:
│       │   ├── Check is not a dir
│       │   └── Check is executable
│       └── Append valid executables to result
└── Return merged list of all executables

Parse() → ([]bool, Arguments)
│
├── Parse --yes / --no / --quit flags
├── Validate only one is set
├── Map to enum: Choice (Yes | No | Quit)
└── Return flag slice and parsed Arguments

CountTrue(flags []bool) → int
│
└── Count how many flags are set to true

```

## Dependencies 
1. Fzf
2. A Unix-type system, like: Macos or Linux, Windows most likely won't work correctly
3. Go >1.24.0
4. Any type of terminal emulator
