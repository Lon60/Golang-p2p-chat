package server

import (
	"Golang-p2p-chat/contact_requests"
	"Golang-p2p-chat/contacts"
	"Golang-p2p-chat/models"
	"Golang-p2p-chat/security"
	"bufio"
	"crypto/rsa"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
)

var PORT = "6000"

func StartServer() {
	// Kommandozeilenargument für die IP-Adresse
	ip := flag.String("ip", "0.0.0.0", "IP-Adresse, auf der der Server lauschen soll")
	port := flag.Int("port", 6000, "Port, auf dem der Server lauschen soll")

	// Parsen der Kommandozeilenargumente
	flag.Parse()

	// Starte den Server auf der angegebenen IP-Adresse und dem Port
	address := *ip + ":" + strconv.Itoa(*port)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("Fehler beim Starten auf %s\n", address)
		return
	}
	defer ln.Close()

	fmt.Printf("Lausche auf %s\n", address)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Fehler beim Akzeptieren der Verbindung:", err)
			continue
		}
		go handleConnection(conn)
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

	contactKeyObject, err := security.ImportPublicKey(contact.PublicKey)
	if err != nil {
		fmt.Println("Fehler beim Importieren des Public Keys:", err)
		return
	}
	contact.KeyObject = contactKeyObject

	fmt.Printf("Kontaktanfrage von %s (%s:%s) wurde akzeptiert.\n", contact.Name, contact.IP, contact.Port)

	// Kontakt hinzufügen
	identifier := contact.IP + ":" + contact.Port
	contacts.AddContact(identifier, contact)
	fmt.Println("Kontakt hinzugefügt.")
}

func handleChatMessage(conn net.Conn, reader *bufio.Reader, senderPublicKey *rsa.PublicKey) {
	// Empfange den Absendernamen
	senderName, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Fehler beim Lesen des Sendernamens:", err)
		return
	}
	senderName = strings.TrimSpace(senderName)

	// Empfange Nachricht und Signatur
	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Fehler beim Lesen der Nachricht:", err)
		return
	}
	message = strings.TrimSpace(message)

	signature, err := reader.ReadBytes('\n')
	if err != nil {
		fmt.Println("Fehler beim Lesen der Signatur:", err)
		return
	}

	// Verifiziere die Signatur mit dem Public Key
	err = security.VerifySignature(message, signature, senderPublicKey)
	if err != nil {
		fmt.Println("Signatur ungültig:", err)
		return
	}

	// Nachricht verarbeiten
	fmt.Printf("Nachricht von %s: %s\n", senderName, message)
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
		// Hier wird angenommen, dass der Absender in den Kontakten gespeichert ist und den Public Key hat
		senderIP, _ := conn.RemoteAddr().(*net.TCPAddr)
		senderContact, exists := contacts.GetContact(senderIP.IP.String() + ":" + strconv.Itoa(senderIP.Port))
		if !exists {
			fmt.Println("Kein Kontakt gefunden, der zu dieser IP passt.")
			return
		}

		// Public Key des Absenders verwenden
		handleChatMessage(conn, reader, senderContact.KeyObject) // Nutze das KeyObject (rsa.PublicKey)
	} else {
		fmt.Println("Unbekannter Nachrichtentyp:", messageType)
	}
}
