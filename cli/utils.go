package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

const (
	ValidatorStatusInactive = 0
	ValidatorStatusActive   = 1
	ValidatorStatusWidrawn  = 2
)

func PanicOnError(message string, err error) {
	if err != nil {
		color.Red(fmt.Sprintf("error: %s\n\n", message))

		info := color.New(color.FgBlack, color.Italic)
		info.Printf(fmt.Sprintf("caused by: %s\n", err))

		os.Exit(1)
	}
}
