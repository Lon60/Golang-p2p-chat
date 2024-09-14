package main

import (
	"Golang-p2p-chat/server"
	"Golang-p2p-chat/ui"
)

func main() {
	go server.StartServer()

	ui.StartUI()
}
