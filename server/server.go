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
	ip := flag.String("ip", "0.0.0.0", "IP address to listen on")
	port := flag.Int("port", 6000, "Port to listen on")

	flag.Parse()

	address := *ip + ":" + strconv.Itoa(*port)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("Error starting on %s\n", address)
		return
	}
	defer ln.Close()

	fmt.Printf("Listening on %s\n", address)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleContactRequest(conn net.Conn, reader *bufio.Reader) {
	requestJSON, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading contact request:", err)
		return
	}

	var request models.ContactRequest
	err = json.Unmarshal([]byte(strings.TrimSpace(requestJSON)), &request)
	if err != nil {
		fmt.Println("Error processing contact request:", err)
		return
	}

	fmt.Printf("Contact request received from %s (%s:%s).\n", request.Name, request.IP, request.Port)
	contact_requests.AddReceivedRequest(request)
}

func handleContactAccepted(conn net.Conn, reader *bufio.Reader) {
	contactJSON, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading contact acceptance:", err)
		return
	}

	var contact models.Contact
	err = json.Unmarshal([]byte(strings.TrimSpace(contactJSON)), &contact)
	if err != nil {
		fmt.Println("Error processing contact acceptance:", err)
		return
	}

	contactKeyObject, err := security.ImportPublicKey(contact.PublicKey)
	if err != nil {
		fmt.Println("Error importing public key:", err)
		return
	}
	contact.KeyObject = contactKeyObject

	fmt.Printf("Contact request from %s (%s:%s) was accepted.\n", contact.Name, contact.IP, contact.Port)

	identifier := contact.IP + ":" + contact.Port
	contacts.AddContact(identifier, contact)
	fmt.Println("Contact added.")
}

func handleChatMessage(conn net.Conn, reader *bufio.Reader, senderPublicKey *rsa.PublicKey) {
	senderName, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading sender name:", err)
		return
	}
	senderName = strings.TrimSpace(senderName)

	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading message:", err)
		return
	}
	message = strings.TrimSpace(message)

	signature, err := reader.ReadBytes('\n')
	if err != nil {
		fmt.Println("Error reading signature:", err)
		return
	}

	err = security.VerifySignature(message, signature, senderPublicKey)
	if err != nil {
		fmt.Println("Invalid signature:", err)
		return
	}

	fmt.Printf("Message from %s: %s\n", senderName, message)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	messageType, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading message type:", err)
		return
	}
	messageType = strings.TrimSpace(messageType)

	if messageType == "CONTACT_REQUEST" {
		handleContactRequest(conn, reader)
	} else if messageType == "CONTACT_ACCEPTED" {
		handleContactAccepted(conn, reader)
	} else if messageType == "CHAT_MESSAGE" {
		senderIP, _ := conn.RemoteAddr().(*net.TCPAddr)
		senderContact, exists := contacts.GetContact(senderIP.IP.String() + ":" + strconv.Itoa(senderIP.Port))
		if !exists {
			fmt.Println("No contact found matching this IP.")
			return
		}

		handleChatMessage(conn, reader, senderContact.KeyObject)
	} else {
		fmt.Println("Unknown message type:", messageType)
	}
}
