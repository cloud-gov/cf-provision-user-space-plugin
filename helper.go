package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Default password length for a password to be generated.
const DefaultPasswordLength = 12

func generatePassword(passwordLength uint32) string {
	// Specify all the acceptable characters to be used in a password.
	acceptableCharacters := []byte{
		// Lowercase
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		// Uppercase
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		// Numbers
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	}
	// Store the length.
	var acceptableCharactersLength = int32(len(acceptableCharacters))

	// Make a slice of strings.
	passwordCharacters := make([]string, passwordLength, passwordLength)

	// Create a seed that makes the numbers nearly random.
	rand.Seed(time.Now().UTC().UnixNano())

	// Iterate through each slot in the slice and store a random accepted character into the slot.
	var index uint32
	for index = 0; index < passwordLength; index++ {
		passwordCharacters[index] = string(acceptableCharacters[rand.Int31n(acceptableCharactersLength)])
	}

	// Return the single joined together string.
	return strings.Join(passwordCharacters, "")
}

func doesFileExist(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

// Helper function to encapsulate what to do when we want to error fatally
func errorPrintln(str string) {
	fmt.Println(str)
	os.Exit(1)
}

func extractUsernameFromEmail(email string) string {
	parts := strings.Split(email, "@")
	// Simple pattern.
	if len(parts) != 2 {
		errorPrintln("Malformed email")
	}
	return parts[0]
}

func interactiveInputValidation() bool {
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(answer)
	switch strings.ToLower(answer) {
	case "y", "yes":
		return true
	case "n", "no":
		errorPrintln("Admin decided to not proceed. Exiting")
		return false
	default:
		errorPrintln("Invalid input. Exiting")
		return false
		// TODO: Currently the cf plugin system prevents from rechecking the input recursively or iteratively.
		// In the meantime, just fail.
		// fmt.Println("Please type 'y' or 'n'")
		// return interactiveInputValidation()
	}

}
