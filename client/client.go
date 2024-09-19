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
	"encoding/base64"
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
	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Printf("error closing connection: %v\n", err)
		}
	}()

	_, err = fmt.Fprintf(conn, "CONTACT_REQUEST\n")
	if err != nil {
		return fmt.Errorf("error sending contact request type: %v", err)
	}

	requestJSON, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("error marshaling contact request: %v", err)
	}

	_, err = fmt.Fprintf(conn, string(requestJSON)+"\n")
	if err != nil {
		return fmt.Errorf("error sending contact request: %v", err)
	}

	return nil
}

func SendContactAccepted(requester models.ContactRequest) error {
	conn, err := net.Dial("tcp", requester.IP+":"+requester.Port)
	if err != nil {
		fmt.Printf("error connecting to %s:%s - %v\n", requester.IP, requester.Port, err)
		return err
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Printf("error closing connection: %v\n", err)
		}
	}()

	_, err = fmt.Fprintf(conn, "CONTACT_ACCEPTED\n")
	if err != nil {
		return fmt.Errorf("error sending contact accepted type: %v", err)
	}

	publicKeyBytes, err := security.ExportPublicKey()
	if err != nil {
		return fmt.Errorf("error exporting public key: %v", err)
	}

	ownContact := models.Contact{
		Name:      config.LocalUserName,
		IP:        getLocalIP(),
		Port:      server.PORT,
		PublicKey: publicKeyBytes,
	}
	contactJSON, err := json.Marshal(ownContact)
	if err != nil {
		return fmt.Errorf("error marshaling own contact: %v", err)
	}

	_, err = fmt.Fprintf(conn, string(contactJSON)+"\n")
	if err != nil {
		return fmt.Errorf("error sending contact accepted: %v", err)
	}

	return nil
}

func SendChatMessage(contact models.Contact, message string) error {
	conn, err := net.Dial("tcp", contact.IP+":"+contact.Port)
	if err != nil {
		messageQueue[contact.Identifier()] = append(messageQueue[contact.Identifier()], message)
		return err
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Printf("error closing connection: %v\n", err)
		}
	}()

	err = sendMessageOverConnection(conn, message, contact.KeyObject)
	if err != nil {
		return fmt.Errorf("error sending chat message: %v", err)
	}

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

func EncryptMessage(message string, recipientPublicKey *rsa.PublicKey) ([]byte, error) {
	encryptedMessage, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, recipientPublicKey, []byte(message), nil)
	if err != nil {
		return nil, fmt.Errorf("error encrypting message: %v", err)
	}
	return encryptedMessage, nil
}

func sendMessageOverConnection(conn net.Conn, message string, recipientPublicKey *rsa.PublicKey) error {

	_, err := fmt.Fprintf(conn, "CHAT_MESSAGE\n")
	if err != nil {
		return fmt.Errorf("error sending messagetype: %v", err)
	}

	_, err = fmt.Fprintf(conn, config.LocalUserName+"\n")
	if err != nil {
		return fmt.Errorf("error sending sender name: %v", err)
	}

	encryptedMessage, err := EncryptMessage(message, recipientPublicKey)
	if err != nil {
		return fmt.Errorf("error encrypted message: %v", err)
	}

	encryptedMessageEncoded := base64.StdEncoding.EncodeToString(encryptedMessage)
	_, err = fmt.Fprintf(conn, encryptedMessageEncoded+"\n")
	if err != nil {
		return fmt.Errorf("error sending encrypted message: %v", err)
	}

	signature, err := signMessage(message)
	if err != nil {
		return fmt.Errorf("error signing message: %v", err)
	}

	signatureEncoded := base64.StdEncoding.EncodeToString(signature)
	_, err = fmt.Fprintf(conn, signatureEncoded+"\n")
	if err != nil {
		return fmt.Errorf("error sending signature: %v", err)
	}

	return nil
}

func signMessage(message string) ([]byte, error) {
	privateKey, err := security.LoadPrivateKey()
	if err != nil {
		return nil, fmt.Errorf("error loading private key: %v", err)
	}

	hash := sha256.New()
	_, err = hash.Write([]byte(message))
	if err != nil {
		return nil, fmt.Errorf("error hashing message: %v", err)
	}
	digest := hash.Sum(nil)

	signature, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, digest, nil)
	if err != nil {
		return nil, fmt.Errorf("error signing message: %v", err)
	}

	return signature, nil
}
