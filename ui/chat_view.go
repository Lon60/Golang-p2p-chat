package ui

import (
	"Golang-p2p-chat/chat_data"
	"Golang-p2p-chat/client"
	"Golang-p2p-chat/config"
	"Golang-p2p-chat/contacts"
	"Golang-p2p-chat/models"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"strings"
)

func viewChatsUI(window fyne.Window) {
	contactsList := contacts.GetContacts()
	if len(contactsList) == 0 {
		dialog.ShowInformation("Keine Chats", "Es sind keine Chats verfügbar.", window)
		return
	}

	var items []fyne.CanvasObject
	for _, contact := range contactsList {
		contactCopy := contact
		item := widget.NewButton(contact.Name, func() {
			startChatUI(window, contactCopy)
		})
		items = append(items, item)
	}

	backButton := widget.NewButton("Zurück", func() {
		showMainMenu(window)
	})

	items = append(items, backButton)
	content := container.NewVBox(items...)
	scroll := container.NewScroll(content)
	window.SetContent(scroll)
}

func startChatUI(window fyne.Window, contact models.Contact) {
	chatHistory := chat_data.GetChatHistory(contact.Identifier())

	messages := widget.NewLabel(strings.Join(*chatHistory, "\n"))
	input := widget.NewEntry()
	input.SetPlaceHolder("Nachricht eingeben")

	sendButton := widget.NewButton("Senden", func() {
		text := strings.TrimSpace(input.Text)
		if text != "" {
			err := client.SendChatMessage(contact, text)
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			chat_data.AppendToChatHistory(contact.Identifier(), "Sie: "+text)
			messages.SetText(strings.Join(*chatHistory, "\n"))
			input.SetText("")
		}
	})

	backButton := widget.NewButton("Zurück", func() {
		viewChatsUI(window)
	})

	content := container.NewVBox(
		messages,
		input,
		sendButton,
		backButton,
	)
	window.SetContent(content)

	go func() {
		for {
			select {
			case newMessage := <-config.MessageChannel:
				chat_data.AppendToChatHistory(contact.Identifier(), newMessage)
				messages.SetText(strings.Join(*chatHistory, "\n"))
			}
		}
	}()
}
