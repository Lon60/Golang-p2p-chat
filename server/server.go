package server

import (
	"Golang-p2p-chat/chat_data"
	"Golang-p2p-chat/config"
	"Golang-p2p-chat/contact_requests"
	"Golang-p2p-chat/contacts"
	"Golang-p2p-chat/models"
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
)

var messageChannel = make(chan string)

var PORT = "6000"
var serverStartMessage string

func StartServer() {
	var ln net.Listener
	var err error

	// Versuchen, den Server auf einem verfügbaren Port zu starten
	portNum, _ := strconv.Atoi(PORT)
	maxPort := portNum + 10 // Wir versuchen Ports von 6000 bis 6010

	for ; portNum <= maxPort; portNum++ {
		PORT = strconv.Itoa(portNum)
		ln, err = net.Listen("tcp", ":"+PORT)
		if err == nil {
			break
		}
		// Keine Ausgabe bei Fehlschlägen
	}

	if err != nil {
		fmt.Println("Fehler beim Starten des Servers:", err)
		return
	}

	// Server läuft erfolgreich
	serverStartMessage = fmt.Sprintf("Server läuft auf Port %s", PORT)
	fmt.Println(serverStartMessage)

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Fehler beim Akzeptieren der Verbindung:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func GetServerStartMessage() string {
	return serverStartMessage
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	messageType, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Fehler beim Lesen des Nachrichtentyps:", err)
		return
	}
	messageType = strings.TrimSpace(messageType)

	if messageType == "CONTACT_REQUEST" {
		handleContactRequest(conn, reader)
	} else if messageType == "CONTACT_ACCEPTED" {
		handleContactAccepted(conn, reader)
	} else if messageType == "CHAT_MESSAGE" {
		handleChatMessage(conn, reader)
	} else {
		fmt.Println("Unbekannter Nachrichtentyp:", messageType)
	}
}

func handleContactRequest(conn net.Conn, reader *bufio.Reader) {
	requestJSON, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Fehler beim Lesen der Kontaktanfrage:", err)
		return
	}

	var request models.ContactRequest
	err = json.Unmarshal([]byte(strings.TrimSpace(requestJSON)), &request)
	if err != nil {
		fmt.Println("Fehler beim Verarbeiten der Kontaktanfrage:", err)
		return
	}

	fmt.Printf("Kontaktanfrage von %s (%s:%s) erhalten.\n", request.Name, request.IP, request.Port)
	contact_requests.AddReceivedRequest(request)
}

func handleContactAccepted(conn net.Conn, reader *bufio.Reader) {
	contactJSON, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Fehler beim Lesen der Kontaktakzeptierung:", err)
		return
	}

	var contact models.Contact
	err = json.Unmarshal([]byte(strings.TrimSpace(contactJSON)), &contact)
	if err != nil {
		fmt.Println("Fehler beim Verarbeiten der Kontaktakzeptierung:", err)
		return
	}

	fmt.Printf("Kontaktanfrage von %s (%s:%s) wurde akzeptiert.\n", contact.Name, contact.IP, contact.Port)

	// Kontakt hinzufügen
	identifier := contact.IP + ":" + contact.Port
	contacts.AddContact(identifier, contact)
	fmt.Println("Kontakt hinzugefügt.")
}

func handleChatMessage(conn net.Conn, reader *bufio.Reader) {
	senderName, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Fehler beim Lesen des Sendernamens:", err)
		return
	}
	senderName = strings.TrimSpace(senderName)

	messageText, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Fehler beim Lesen der Nachricht:", err)
		return
	}
	messageText = strings.TrimSpace(messageText)

	fullMessage := fmt.Sprintf("%s: %s", senderName, messageText)
	chat_data.AppendToChatHistory(senderName, fullMessage)

	// Sende die Nachricht an das UI über den Channel
	config.MessageChannel <- fullMessage

	fmt.Printf("\nNeue Nachricht von %s erhalten: %s\n", senderName, messageText)
}
