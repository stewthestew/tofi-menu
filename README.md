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
Tofi-menu uses these different environment variables:

1. LAUNCH_APP: The path to the fzf binary. (Optional, if not set then it will just execute fzf normally).
2. DEBUG: If this is set, then less important errors will be handled (Read code for more info) (Optional).
3. TOFI_APPS: This acts like the PATH variable, this is where tofi will get the executables from.
4. SILENT: This silences the output of the program, it will not print anything to the terminal.

If you want it to list every app then do 
```bash
export TOFI_APPS=$PATH
```

You can also change TOFI_APPS to point at specific directories so you can launch custom scripts, etc.

## Dependencies 
1. Fzf
2. A Unix-type system, like: Macos or Linux, Windows most likely won't work correctly
3. Go >1.24.0
4. Any type of terminal emulator
