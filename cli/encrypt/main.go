package main

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/hlandau/passlib.v1"
	"syscall"
)

func main() {
	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		panic(err)
	}
	hash, err := passlib.Hash(string(bytePassword))
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n%s\n", hash)
}
