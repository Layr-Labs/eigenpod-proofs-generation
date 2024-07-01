package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func PanicOnError(message string, err error) {
	if err != nil {
		color.Red(fmt.Sprintf("error: %s\n\n", message))

		info := color.New(color.FgBlack, color.Italic)
		info.Printf(fmt.Sprintf("caused by: %s\n", err))

		os.Exit(1)
	}
}
