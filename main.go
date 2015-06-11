package main

import (
	"flag"
	"fmt"
	"github.com/cloudfoundry/cli/plugin"
	"os"
	"runtime"
)

// ProvisionUserPlugin is the struct that implements the interface defined by the cf core CLI.
// The interface can be found at: github.com/cloudfoundry/cli/plugin/plugin.go
type ProvisionUserPlugin struct{}

// ProvisionSingleUser will create the user, org and space for a single user.
func (c *ProvisionUserPlugin) ProvisionSingleUser(userdata *userData, cliConnection *plugin.CliConnection) {
}

func (c *ProvisionUserPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	// Check for compatible OS.
	if runtime.GOOS == "windows" {
		errorPrintln("Detected incompatible OS: Windows. Exiting...")
	}

	/*
		PARSE INPUT
	*/
	flagSet := flag.NewFlagSet("provision-user-space", flag.ContinueOnError)
	emailFlag := flagSet.String("email", "", "The specified e-mail address of the account to be created")
	orgFlag := flagSet.String("org", "", "The specified org of the account to be created")
	_ = flagSet.Parse(args[1:])
	if len(*emailFlag) < 1 {
		errorPrintln("No email flag given. Please run with --help for usage.")
	}
	if len(*orgFlag) < 1 {
		errorPrintln("No org flag given. Please run with --help for usage.")
	}
	userdata := userData{email: *emailFlag, username: extractUsernameFromEmail(*emailFlag), org: *orgFlag}

	/*
		SETUP
	*/
	// Check if logged in.
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

	downloadFugu()

	/*
		SUMMARY
	*/
	// Give the admin a summary before the actions are applied.
	fmt.Println("Your username: " + username)
	userdata.printIncompleteUserData()


	// Validate that indeed they want to proceed.
	var interactive bool = true
	if interactive {
		// var answer string
		fmt.Println("Is this correct? [y/n]")
		var proceed bool = interactiveInputValidation()
		if !proceed {
			errorPrintln("Admin decided to not proceed. Exiting")
		}
	}
	// fmt.Println("Made it out")

	/*
		EXECUTE
	*/
	// TODO

	os.Exit(0)
}

func (c *ProvisionUserPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "Provision-User-Space-Plugin",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 0,
		},
		Commands: []plugin.Command{
			plugin.Command{
				Name:     "provision-user-space",
				HelpText: "This plugin creates the specified user and org and a personal space. ",
				UsageDetails: plugin.Usage{
					Usage: "cf provision-user-space [-email=<username@domain>] [-org=<org>]",
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(ProvisionUserPlugin))
}
