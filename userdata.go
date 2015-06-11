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

func (user *userData) printPartialUserData() {
	fmt.Printf("email: (%s) username: (%s) org: (%s)\n", user.email, user.username, user.org)
}

func (user *userData) printUserData() {
	fmt.Printf("email: (%s) username: (%s) org: (%s) fugu-url: (%s) \n", user.email, user.username, user.org, user.fuguURL)
}
