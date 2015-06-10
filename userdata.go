package main

import (
	"fmt"
)

type UserData struct {
	email    string
	password string
	fuguURL  string
}

func (user *UserData) getEmail() string {
	return "email: " + user.email
}

func (user *UserData) printPublicData() {
	fmt.Printf("%s fuguURL: %s", user.getEmail(), user.fuguURL)
}
