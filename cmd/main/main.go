package main

import (
	lib "tofi/internal/tofiLib"
)

func main() {
	o := lib.TofiOptionsInit()
	e := lib.EnvVarsInit(o.Silent)
	lib.Tofi(o, e)
}
