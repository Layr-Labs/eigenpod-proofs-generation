package main

import "fmt"

func PanicOnError(message string, err error) {
	if err != nil {
		panic(fmt.Errorf("%s\n\ncaused by: %w", message, err))
	}
}
