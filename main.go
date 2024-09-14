package main

import (
	"Golang-p2p-chat/security"
	"Golang-p2p-chat/server"
	"Golang-p2p-chat/ui"
	"fmt"
)

func main() {
	err := security.GenerateKeyPairIfNotExists(2048) // 2048 Bits Schlüssellänge
	if err != nil {
		fmt.Println("Fehler beim Generieren des Schlüsselpaares:", err)
		return
	}

	go server.StartServer()

	ui.StartUI()
}
