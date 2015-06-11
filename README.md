# cf-provision-user-space

This is a [plugin](https://github.com/cloudfoundry/cli/tree/master/plugin_examples) for [Cloud Foundry CLI](https://github.com/cloudfoundry/cli). This plugin effectively creates a new cf user, an organization (if specified) and a personal space within the sandbox organization. At the end of the script, it will print out.

	email: (<the specified e-mail>) username: (<username section from the e-mail>) org: (<the optional org>) fugu-url: (<the time limited URL to where to find the temporary password>)

This plugin is a port of this [shell script](https://github.com/18F/cloud-foundry-scripts/blob/674a511662490165e629d77fb4e9dda28837b27a/cf-create-user.sh).

## Pre-Requisites

### CloudFoundry CLI

To install the CLI, follow the instructions [here](https://github.com/cloudfoundry/cli).

Once installed, make sure you are logged into your API server.

	$ cf login

Also, make sure your API server has a sandbox organization. You can check via

	$ cf orgs

### Go

Install [Go](https://golang.org/).

On Mac OS X with brew

	$ brew install go


## Installation
To install the plugin:

	$ go get github.com/18F/cf-provision-user-space-plugin
	$ cf install-plugin $GOBIN/cf-provision-user-space-plugin

## Usage
To use the plugin:

	$ cf provision-user-space [-email=<username@domain>] [-org=<org> (optional)]


## Removal
To remove the plugin:

	$ cf uninstall-plugin "Provision-User-Space-Plugin"
