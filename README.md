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
