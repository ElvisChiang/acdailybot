package main

import (
	"errors"
	"fmt"
	"strings"
)

// Command main entry for parsing command
func Command(who string, msg string) (err error) {
	err = errors.New("execuse me?")

	result := strings.Fields(msg)
	if len(result) == 0 {
		return
	}
	command := result[0]
	lowerCmd := strings.ToLower(command)
	// lowerCmd = strings.Replace(lowerCmd, "@walkrbot", "", -1)
	fmt.Printf("command [%s] %s\n", who, lowerCmd)

	return
}

// Highlight my data /hl
func Highlight() (ok bool) {
	ok = true
	return
}

// Remove my data /remove
func Remove() (ok bool) {
	ok = true

	return
}
