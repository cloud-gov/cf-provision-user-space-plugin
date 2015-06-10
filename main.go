package main

import (
	"fmt"
	"github.com/cloudfoundry/cli/plugin"
	"os"
	"runtime"
)

// ProvisionUserPlugin is the struct that implements the interface defined by the cf core CLI.
// The interface can be found at: github.com/cloudfoundry/cli/plugin/plugin.go
type ProvisionUserPlugin struct{}

// ProvisionSingleUser will create the user, org and space for a single user.
func (c *ProvisionUserPlugin) ProvisionSingleUser(userdata *UserData, cliConnection *plugin.CliConnection) {
}



func (c *ProvisionUserPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	// Check for compatible OS.
	if runtime.GOOS == "windows" {
		errorPrintln("Detected incompatible OS: Windows. Exiting...")
	}

	// Check if logged in.
	fmt.Printf("Found the following args %v\n", args)
	loggedIn, err := cliConnection.IsLoggedIn()
	if err != nil {
		errorPrintln("Unable to connect.")
	}
	if !loggedIn {
		errorPrintln("Not logged in. Please run 'cf login'")
	}

	// Get the administrator's username.
	username, err := cliConnection.Username()
	// TODO. Check err in case we lose connection to the server intermittently.
	if len(username) == 0 {
		errorPrintln("No username found.")
	}

	// Give the admin a summary before the actions are applied.
	fmt.Println("Your username: " + username)

	//newUserEmail := "test"
	// newPassword := generatePassword()

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
			errorPrintln("Admin decided to not proceed. Exiting")
		}
	}
	fmt.Println("Made it out")
	fmt.Println("JaMES we are lucky " + downloadFuguAndUploadPassword("test"))
	os.Exit(0)
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
