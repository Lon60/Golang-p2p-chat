package client

import (
	"Golang-p2p-chat/config"
	"Golang-p2p-chat/models"
	"Golang-p2p-chat/server"
	"encoding/json"
	"fmt"
	"net"
)

func SendContactRequest(ip, port string, request models.ContactRequest) error {
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Senden des Nachrichtentyps
	fmt.Fprintf(conn, "CONTACT_REQUEST\n")

	// Senden der Kontaktanfrage als JSON
	requestJSON, _ := json.Marshal(request)
	fmt.Fprintf(conn, string(requestJSON)+"\n")

	return nil
}

func SendContactAccepted(requester models.ContactRequest) error {
	conn, err := net.Dial("tcp", requester.IP+":"+requester.Port)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Senden des Nachrichtentyps
	fmt.Fprintf(conn, "CONTACT_ACCEPTED\n")

	// Senden der eigenen Kontaktdaten als JSON
	ownContact := models.Contact{
		Name: config.LocalUserName,
		IP:   getLocalIP(),
		Port: server.PORT,
	}
	contactJSON, _ := json.Marshal(ownContact)
	fmt.Fprintf(conn, string(contactJSON)+"\n")

	return nil
}

func SendChatMessage(contact models.Contact, message string) error {
	conn, err := net.Dial("tcp", contact.IP+":"+contact.Port)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Senden des Nachrichtentyps
	fmt.Fprintf(conn, "CHAT_MESSAGE\n")

	// Senden des eigenen Namens
	fmt.Fprintf(conn, config.LocalUserName+"\n")

	// Senden der Nachricht
	fmt.Fprintf(conn, message+"\n")

	return nil
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String()
		}
	}
	return "127.0.0.1"
}
