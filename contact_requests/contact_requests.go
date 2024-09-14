package contact_requests

import (
	"Golang-p2p-chat/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var receivedRequests = make(map[string]models.ContactRequest)
var requestsFile = "received_requests.json"

func AddReceivedRequest(request models.ContactRequest) {
	LoadRequestsFromFile()
	identifier := request.IP + ":" + request.Port
	receivedRequests[identifier] = request
	SaveRequestsToFile()
}

func GetReceivedRequests() []models.ContactRequest {
	LoadRequestsFromFile()
	requests := []models.ContactRequest{}
	for _, req := range receivedRequests {
		requests = append(requests, req)
	}
	return requests
}

func RemoveReceivedRequestByIdentifier(identifier string) {
	LoadRequestsFromFile()
	delete(receivedRequests, identifier)
	SaveRequestsToFile()
}

func SaveRequestsToFile() {
	file, err := json.MarshalIndent(receivedRequests, "", "  ")
	if err != nil {
		fmt.Println("Fehler beim Speichern der Kontaktanfragen:", err)
		return
	}
	err = ioutil.WriteFile(requestsFile, file, 0644)
	if err != nil {
		fmt.Println("Fehler beim Schreiben der Kontaktanfragen-Datei:", err)
	}
}

func LoadRequestsFromFile() {
	if _, err := os.Stat(requestsFile); os.IsNotExist(err) {
		return
	}
	file, err := ioutil.ReadFile(requestsFile)
	if err != nil {
		fmt.Println("Fehler beim Lesen der Kontaktanfragen-Datei:", err)
		return
	}
	err = json.Unmarshal(file, &receivedRequests)
	if err != nil {
		fmt.Println("Fehler beim Laden der Kontaktanfragen:", err)
	}
}
