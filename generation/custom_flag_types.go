package main

import (
	"fmt"
	"strconv"
	"strings"
)

type IntSlice struct {
	slice []uint64
}

func (i *IntSlice) String() string {
	return fmt.Sprintf("%v", i.slice)
}

func (i *IntSlice) GetSlice() []uint64 {
	return i.slice
}

func (i *IntSlice) Set(value string) error {
	// Split the string by commas if user passes a comma-separated list.
	vals := strings.Split(value, ",")
	for _, val := range vals {
		trimmedVal := strings.TrimSpace(val)
		if trimmedVal == "" {
			continue
		}
		intVal, err := strconv.ParseUint(trimmedVal, 10, 64)
		if err != nil {
			return err
		}
		i.slice = append(i.slice, intVal)
	}
	return nil
}
