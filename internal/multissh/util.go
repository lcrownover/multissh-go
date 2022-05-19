package multissh

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
)

// getSecurePassword : get password input securely
func getSecurePassword(prompt string) (string, error) {
	fmt.Printf(prompt)
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Printf("\n")
	if err != nil {
		logger("config").Error("failed to get secure password")
		return "", err
	}
	password := string(bytePassword)
	return password, nil
}
