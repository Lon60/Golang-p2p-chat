package client

import (
	"Golang-p2p-chat/config"
	"Golang-p2p-chat/models"
	"Golang-p2p-chat/security"
	"Golang-p2p-chat/server"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net"
)

var messageQueue = make(map[string][]string)

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
		fmt.Printf("Fehler bei der Verbindung zu %s:%s - %v\n", requester.IP, requester.Port, err)
		return err
	}
	defer conn.Close()

	// Senden des Nachrichtentyps
	fmt.Fprintf(conn, "CONTACT_ACCEPTED\n")

	// Lade den eigenen Public Key
	publicKeyBytes, err := security.ExportPublicKey()
	if err != nil {
		return err
	}

	// Senden der eigenen Kontaktdaten und Public Key als JSON
	ownContact := models.Contact{
		Name:      config.LocalUserName,
		IP:        getLocalIP(),
		Port:      server.PORT,
		PublicKey: publicKeyBytes, // Sende Public Key als []byte
	}
	contactJSON, _ := json.Marshal(ownContact)
	fmt.Fprintf(conn, string(contactJSON)+"\n")

	return nil
}

func SendChatMessage(contact models.Contact, message string) error {
	// Verbindung aufbauen
	conn, err := net.Dial("tcp", contact.IP+":"+contact.Port)
	if err != nil {
		// Wenn der Kontakt offline ist, Nachricht in die Warteschlange legen
		messageQueue[contact.Identifier()] = append(messageQueue[contact.Identifier()], message)
		return err
	}
	defer conn.Close()

	// Nachricht senden
	return sendMessageOverConnection(conn, message)
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

func EncryptMessage(message string, recipientPublicKey *rsa.PublicKey) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, recipientPublicKey, []byte(message), nil)
}

func sendMessageOverConnection(conn net.Conn, message string) error {
	// Absendername senden
	_, err := fmt.Fprintf(conn, config.LocalUserName+"\n")
	if err != nil {
		return fmt.Errorf("Fehler beim Senden des Absendernamens: %v", err)
	}

	// Nachricht senden
	_, err = fmt.Fprintf(conn, message+"\n")
	if err != nil {
		return fmt.Errorf("Fehler beim Senden der Nachricht: %v", err)
	}

	// Nachricht signieren
	signature, err := signMessage(message)
	if err != nil {
		return fmt.Errorf("Fehler beim Signieren der Nachricht: %v", err)
	}

	// Signatur senden
	_, err = conn.Write(signature)
	if err != nil {
		return fmt.Errorf("Fehler beim Senden der Signatur: %v", err)
	}

	return nil
}

// signMessage signiert eine Nachricht mit dem Private Key des Benutzers.
func signMessage(message string) ([]byte, error) {
	privateKey, err := security.LoadPrivateKey() // Lade den Private Key des Benutzers
	if err != nil {
		return nil, err
	}

	hash := sha256.New()
	hash.Write([]byte(message))
	digest := hash.Sum(nil)

	// Signiere die Nachricht
	signature, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, digest, nil)
	if err != nil {
		return nil, fmt.Errorf("Fehler beim Signieren der Nachricht: %v", err)
	}

	return signature, nil
}
