package contact_requests

import (
	"Golang-p2p-chat/models"
	"encoding/json"
	"fmt"
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
		fmt.Println("Error saving contact requests:", err)
		return
	}
	err = os.WriteFile(requestsFile, file, 0644)
	if err != nil {
		fmt.Println("Error writing contact requests file:", err)
	}
}

func LoadRequestsFromFile() {
	if _, err := os.Stat(requestsFile); os.IsNotExist(err) {
		return
	}
	file, err := os.ReadFile(requestsFile)
	if err != nil {
		fmt.Println("Error reading contact requests file:", err)
		return
	}
	err = json.Unmarshal(file, &receivedRequests)
	if err != nil {
		fmt.Println("Error loading contact requests:", err)
	}
}
