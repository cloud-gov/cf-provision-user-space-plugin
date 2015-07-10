package main

import (
	"flag"
	"fmt"
	"github.com/cloudfoundry/cli/plugin"
	"os"
	"runtime"
	"strings"
)

// ProvisionUserPlugin is the struct that implements the interface defined by the cf core CLI.
// The interface can be found at: github.com/cloudfoundry/cli/plugin/plugin.go
type ProvisionUserPlugin struct{}

// ProvisionSingleUser will create the user, org and space for a single user.
func (c *ProvisionUserPlugin) ProvisionSingleUser(userdata *userData, cliConnection plugin.CliConnection) {
	// Create password.
	userdata.password = generatePassword(DefaultPasswordLength)

	// Create user.
	output, err := cliConnection.CliCommandWithoutTerminalOutput("create-user", userdata.email, userdata.password)

	if err != nil {
		errorPrintln("Unable to create user. Please check if you have admin access. Error: " + err.Error())
	} else {
		userAlreadyExists := false
		// Check if output says the user already exists.
		for _, line := range output {
			if strings.Contains(line, "already exists") {
				fmt.Println("User " + userdata.email + " already exists. Will not attempt to create user.")
				userAlreadyExists = true
				break
			}
		}

		// Upload password with fugu if the user was just created.
		if !userAlreadyExists {
			fmt.Println("Uploading with fugu")
			userdata.fuguURL = uploadPasswordToFugu(userdata.password)
		}
	}

	// If the spaces doesn't exist, create it.
	foundUserSpace := false
	spaces, err := cliConnection.CliCommandWithoutTerminalOutput("spaces")
	for _, space := range spaces {
		if space == userdata.username {
			foundUserSpace = true
			break
		}
	}
	if foundUserSpace {
		fmt.Println("Space '" + userdata.username + "' already exists")
	} else {
		_, err = cliConnection.CliCommandWithoutTerminalOutput("create-space", userdata.username)
		if err != nil {
			errorPrintln("Unable to create space '" + userdata.username + "' Error: " + err.Error())
		}
	}
	// Set user permissions.
	_, err = cliConnection.CliCommandWithoutTerminalOutput("set-space-role", userdata.email, "sandbox", userdata.username, "SpaceDeveloper")
	_, err = cliConnection.CliCommandWithoutTerminalOutput("set-space-role", userdata.email, "sandbox", userdata.username, "SpaceManager")

	// Create the org if supplied and give the user permissions.
	if len(userdata.org) > 0 {
		foundOrg := false
		// Get the orgs.
		orgs, err := cliConnection.CliCommandWithoutTerminalOutput("orgs")
		// Check to see if the org already exists.
		for _, org := range orgs {
			if org == userdata.org {
				foundOrg = true
				break
			}
		}
		if foundOrg {
			fmt.Println("Org '" + userdata.org + "' already exists")
		} else {
			_, err = cliConnection.CliCommandWithoutTerminalOutput("create-org", userdata.org)
			if err != nil {
				errorPrintln("Unable to create org '" + userdata.org + "' Error: " + err.Error())
			}
		}

		// Set org roles.
		_, err = cliConnection.CliCommandWithoutTerminalOutput("set-org-role", userdata.email, userdata.org, "OrgManager")
		// Since the typical expectation is that being OrgManager confers access to the contained
		// spaces as well, but doesn't we'll go ahead and add those permissions.
		cliConnection.CliCommandWithoutTerminalOutput("target", "-o", userdata.org)
		spaces, _ = cliConnection.CliCommandWithoutTerminalOutput("spaces")
		for _, space := range spaces {
			_, err = cliConnection.CliCommandWithoutTerminalOutput("set-space-role", userdata.email, userdata.org, space, "SpaceDeveloper")
			_, err = cliConnection.CliCommandWithoutTerminalOutput("set-space-role", userdata.email, userdata.org, space, "SpaceManager")
		}
	}

	// Finish with a print.
	userdata.printUserData()
}

func (c *ProvisionUserPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	// Check for compatible OS.
	if runtime.GOOS == "windows" {
		errorPrintln("Detected incompatible OS: Windows. Exiting...")
	}

	/*
		PARSE INPUT
	*/
	// Arg 1: provision-user-space
	// Args 2-N: args
	if len(args) < 2 {
		errorPrintln("Not enough args")
	}

	// Get the e-mail address.
	email := args[1]

	// Setup the optional flags.
	flagSet := flag.NewFlagSet("provision-user-space", flag.ContinueOnError)
	orgFlag := flagSet.String("org", "", "The specified org of the account to be created")
	// Only look at the other arguments if we have more.
	if len(args) > 2 {
		_ = flagSet.Parse(args[2:])
		if len(*orgFlag) < 1 {
			// TODO warn but don't fail.
		}
	}

	userdata := userData{email: email, username: extractUsernameFromEmail(email), org: *orgFlag}

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

	_, err = cliConnection.CliCommandWithoutTerminalOutput("target", "-o", "sandbox")
	if err != nil {
		errorPrintln("Unable to target the sandbox org")
	}
	downloadFugu()

	/*
		SUMMARY
	*/
	// Give the admin a summary before the actions are applied.
	fmt.Println("Your username: " + username)
	userdata.printPartialUserData()

	// Validate that indeed they want to proceed.
	var interactive = true
	if interactive {
		fmt.Println("Is this correct? [y/n]")
		var proceed = interactiveInputValidation()
		if !proceed {
			errorPrintln("Admin decided to not proceed. Exiting")
		}
	}

	/*
		EXECUTE
	*/
	c.ProvisionSingleUser(&userdata, cliConnection)
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
					Usage: "cf provision-user-space <username@domain> [-org=<org> (optional)]",
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(ProvisionUserPlugin))
}
