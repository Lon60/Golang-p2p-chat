package main

import (
	"Golang-p2p-chat/security"
	"Golang-p2p-chat/server"
	"Golang-p2p-chat/ui"
	"fmt"
)

func main() {
	err := security.GenerateKeyPairIfNotExists(2048)
	if err != nil {
		fmt.Println("Error generating key pair:", err)
		return
	}

	go server.StartServer()

	ui.StartUI()
}
