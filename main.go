package main

import (
	"bufio"
	"fmt"
	"github.com/cloudfoundry/cli/plugin"
	"os"
	"strings"
)

// ProvisionUserPlugin is the struct that implements the interface defined by the cf core CLI.
// The interface can be found at: github.com/cloudfoundry/cli/plugin/plugin.go
type ProvisionUserPlugin struct{}

// ProvisionSingleUser will create the user, org and space for a single user.
func (c *ProvisionUserPlugin) ProvisionSingleUser() {
}

// Helper function to encapsulate what to do when we want to error fatally
func errorPrintln(str string) {
	fmt.Println(str)
}

func interactiveInputValidation() (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(answer)
	switch strings.ToLower(answer) {
	case "y", "yes":
		return true
		// proceed = true
		// break InteractiveInputValidation
	case "n", "no":
		return false
		// proceed = false
		// break InteractiveInputValidation
	default:
		return false
		// TODO: Currently the cf plugin system prevents from rechecking the input recurisively or iteratively.
		// In the meantime, just fail.
		// fmt.Println("Please type 'y' or 'n'")
		// return interactiveInputValidation()
	}

}
func (c *ProvisionUserPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	// output, err := cliConnection.CliCommandWithoutTerminalOutput("apps")
	loggedIn, err := cliConnection.IsLoggedIn()
	if err != nil {
		errorPrintln("Unable to connect.")
		return
	}
	if !loggedIn {
		errorPrintln("Not logged in. Please run 'cf login'")
		return
	}

	username, err := cliConnection.Username()
	// TODO. Check err in case we lose connection to the server intermittenly.
	if len(username) == 0 {
		errorPrintln("No username found.")
		return
	}

	// Give the admin a summary before the actions are applied.
	fmt.Println("Your username: " + username)

	// Validate that indeed they want to proceed.
	var interactive bool = true
	if interactive {
		// var answer string
		fmt.Println("Is this correct? [y/n]")
		var proceed bool = interactiveInputValidation()
		/*
				InteractiveInputValidation:
			reader := bufio.NewReader(os.Stdin)
					for {
						answer, _ := reader.ReadString('\n')
						answer = strings.TrimSpace(answer)
						fmt.Print("found (" + answer +")")
						switch strings.ToLower(answer) {
						case "y", "yes":
							proceed = true
							break InteractiveInputValidation
						case "n", "no":
							proceed = false
							break InteractiveInputValidation
						}
						fmt.Println("Please type 'y' or 'n'")
					}

		*/
		if !proceed {
			fmt.Println("Admin decided to not proceed. Exiting")
		}
	}
	fmt.Println("Made it out")
}

func (c *ProvisionUserPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "Provision-New-User-Plugin",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 0,
		},
		Commands: []plugin.Command{
			plugin.Command{
				Name:     "provision-user-workspace",
				HelpText: "some text for now",
			},
		},
	}
}

func main() {
	plugin.Start(new(ProvisionUserPlugin))
}
