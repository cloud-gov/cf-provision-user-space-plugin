package main

import (
	"fmt"
)

type userData struct {
	email    string
	username string
	org      string
	password string
	fuguURL  string
}

func (user *userData) printIncompleteUserData() {
	fmt.Printf("email: %s username: %s org: %s\n", user.email, user.username, user.org)
}
