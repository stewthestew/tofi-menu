package utils

import (
	"github.com/charmbracelet/log"
)

func Log(level log.Level, silent bool, msg any, keyvals ...any) {
	if !silent {
		log.Log(level, msg, keyvals)
	}
}
