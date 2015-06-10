package main

import (
	"path/filepath"
	"fmt"
	"os"
	"io"
	"path"
	"net/http"
	"os/exec"
)

func downloadFuguAndUploadPassword(password string) string {
	if !doesFuguExist() {
		fmt.Println("Fugu does not exist. Attempting to download")
		downloadFugu()
	}
	return uploadPasswordToFugu(password)
}

func uploadPasswordToFugu(password string) (string) {
	// TODO If password could ever accept whitespaces, we need to surround with single quotes
	// so it remains as one argument.
	output, err := exec.Command(getFuguPath(), password).Output()
	if err != nil {
		errorPrintln("Cannot execute fugu. Please check internet connection and if binary exist at: " +getFuguPath() + "Error: " + err.Error())
	}
	return string(output[:])
}

// Downloads the Fugu binary which will upload the password.
func downloadFugu() {
	// Create folder(s) if needed.
	fuguBinPath := getFuguPath()
	fuguFolders := filepath.Dir(fuguBinPath)
	if (!doesFileExist(fuguFolders)) {
		fmt.Println("Making folders at path: " + fuguFolders)
		if os.MkdirAll(fuguFolders, 0700) != nil {
			errorPrintln("Unable to make folder for fugu at path: (" + fuguFolders + ")")
		}
	}

	// Create file.
	file, err := os.Create(fuguBinPath)
	if err != nil {
		errorPrintln("Unable to create fugu binary file. error: " + err.Error())
	}
	defer file.Close()

	// Download into file.
	fuguURL := "https://raw.githubusercontent.com/jgrevich/fugacious/feature/14_fugacious-CLI/bin/fugu"
	response, err := http.Get(fuguURL)
	if err != nil {
		os.Remove(fuguBinPath)
		errorPrintln("Unable to download fugu from: " + fuguURL)
	}
	defer response.Body.Close()

	_, err = io.Copy(file, response.Body)

	// Change permissions to execute.
	file.Chmod(0700)
}

func getFuguPath() string {
	// TODO allow to override the path from command line
	homePath := os.Getenv("HOME")

	// Quick error check to see if we were able to find the HOME variable.
	if len(homePath) < 1 {
		errorPrintln("Can not find HOME environment variable")
	}

	// Return the joined path.
	return path.Join(homePath, ".cf", "bin", "fugu")
}

func doesFuguExist() bool {
	// TODO Check PATH.

	// Check $HOME/.cf/bin/fugu
	return doesFileExist(getFuguPath())
}
