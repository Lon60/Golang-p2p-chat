package ui

import (
	"Golang-p2p-chat/contacts"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"strings"
)

func editContactNameUI(window fyne.Window) {
	contactsList := contacts.GetContacts()
	if len(contactsList) == 0 {
		dialog.ShowInformation("No Contacts", "No contacts are available.", window)
		return
	}

	var items []fyne.CanvasObject
	for identifier, contact := range contactsList {
		identifierCopy := identifier
		item := widget.NewButton(fmt.Sprintf("%s (current: %s)", identifier, contact.Name), func() {
			entry := widget.NewEntry()
			entry.SetPlaceHolder("New Name")

			form := dialog.NewForm("Edit Contact Name", "Save", "Cancel", []*widget.FormItem{
				widget.NewFormItem("New Name", entry),
			}, func(b bool) {
				if b {
					newName := strings.TrimSpace(entry.Text)
					if newName != "" {
						contacts.UpdateContactName(identifierCopy, newName)
						dialog.ShowInformation("Success", "Contact name updated.", window)
					} else {
						dialog.ShowInformation("Error", "Name cannot be empty.", window)
					}
				}
			}, window)
			form.Show()
		})
		items = append(items, item)
	}

	backButton := widget.NewButton("Back", func() {
		showMainMenu(window)
	})

	items = append(items, backButton)
	content := container.NewVBox(items...)
	scroll := container.NewScroll(content)
	window.SetContent(scroll)
}
