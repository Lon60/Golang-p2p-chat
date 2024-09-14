package contacts

import (
	"Golang-p2p-chat/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var contacts = make(map[string]models.Contact)
var contactsFile = "contacts.json"

func AddContact(identifier string, contact models.Contact) {
	contacts[identifier] = contact
	SaveContactsToFile()
}

func GetContact(identifier string) (models.Contact, bool) {
	contact, exists := contacts[identifier]
	return contact, exists
}

func GetContacts() map[string]models.Contact {
	LoadContactsFromFile()
	return contacts
}

func UpdateContactName(identifier, newName string) {
	if contact, exists := contacts[identifier]; exists {
		contact.Name = newName
		contacts[identifier] = contact
		SaveContactsToFile()
	}
}

func SaveContactsToFile() {
	file, err := json.MarshalIndent(contacts, "", "  ")
	if err != nil {
		fmt.Println("Fehler beim Speichern der Kontakte:", err)
		return
	}
	err = ioutil.WriteFile(contactsFile, file, 0644)
	if err != nil {
		fmt.Println("Fehler beim Schreiben der Kontakte-Datei:", err)
	}
}

func LoadContactsFromFile() {
	if _, err := os.Stat(contactsFile); os.IsNotExist(err) {
		return
	}
	file, err := ioutil.ReadFile(contactsFile)
	if err != nil {
		fmt.Println("Fehler beim Lesen der Kontakte-Datei:", err)
		return
	}
	err = json.Unmarshal(file, &contacts)
	if err != nil {
		fmt.Println("Fehler beim Laden der Kontakte:", err)
	}
}
