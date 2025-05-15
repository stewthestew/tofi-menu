package main

import (
	"tofi/internal/tofiLib"
)

func main() {
	opts := tofiLib.TofiOptionsInit()
	e := tofiLib.EnvVarsInit(opts.Silent)
	tofiLib.Tofi(opts, e)
}
