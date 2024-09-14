package chat_data

import (
	"sync"
)

var (
	Chats = make(map[string]*[]string)
	mutex = &sync.Mutex{}
)

func GetChatHistory(identifier string) *[]string {
	mutex.Lock()
	defer mutex.Unlock()

	if _, exists := Chats[identifier]; !exists {
		var chatHistory []string
		Chats[identifier] = &chatHistory
	}
	return Chats[identifier]
}

func AppendToChatHistory(identifier, message string) {
	mutex.Lock()
	defer mutex.Unlock()

	if _, exists := Chats[identifier]; !exists {
		var chatHistory []string
		Chats[identifier] = &chatHistory
	}
	*Chats[identifier] = append(*Chats[identifier], message)
}
