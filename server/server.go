package server

import (
	"Golang-p2p-chat/chat_data"
	"Golang-p2p-chat/config"
	"Golang-p2p-chat/contact_requests"
	"Golang-p2p-chat/contacts"
	"Golang-p2p-chat/models"
	"Golang-p2p-chat/security"
	"bufio"
	"crypto/rsa"
	"encoding/base64"
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

func handleContactRequest(reader *bufio.Reader) {
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

func handleContactAccepted(reader *bufio.Reader) {
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

func handleChatMessage(reader *bufio.Reader, senderPublicKey *rsa.PublicKey, senderName string) {
	senderName, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading sender name:", err)
		return
	}
	senderName = strings.TrimSpace(senderName)

	encryptedMessageEncoded, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("error reading message:", err)
		return
	}

	encryptedMessage, err := base64.StdEncoding.DecodeString(strings.TrimSpace(encryptedMessageEncoded))
	if err != nil {
		fmt.Println("error decode message:", err)
		return
	}

	message, err := security.DecryptMessage(encryptedMessage)
	if err != nil {
		fmt.Println("error decrypting message:", err)
		return
	}

	signatureEncoded, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("error reading signature:", err)
		return
	}

	signature, err := base64.StdEncoding.DecodeString(strings.TrimSpace(signatureEncoded))
	if err != nil {
		fmt.Println("error decoding signature:", err)
		return
	}

	err = security.VerifySignature(message, signature, senderPublicKey)
	if err != nil {
		fmt.Println("Invalid signature:", err)
		return
	}

	fmt.Printf("Message from %s: %s\n", senderName, message)
	chat_data.AppendToChatHistory(senderName, senderName+": "+message)
	config.MessageChannel <- senderName + ": " + message
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
		handleContactRequest(reader)
	} else if messageType == "CONTACT_ACCEPTED" {
		handleContactAccepted(reader)
	} else if messageType == "CHAT_MESSAGE" {
		senderName, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading sender name:", err)
			return
		}
		senderName = strings.TrimSpace(senderName)

		senderContact, exists := contacts.GetContactByName(senderName)
		if !exists {
			fmt.Println("No contact found matching this name.")
			return
		}

		handleChatMessage(reader, senderContact.KeyObject, senderName)
	} else {
		fmt.Println("Unknown message type:", messageType)
	}
}
