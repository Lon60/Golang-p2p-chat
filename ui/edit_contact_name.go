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
		dialog.ShowInformation("Keine Kontakte", "Es sind keine Kontakte verfügbar.", window)
		return
	}

	var items []fyne.CanvasObject
	for identifier, contact := range contactsList {
		identifierCopy := identifier
		item := widget.NewButton(fmt.Sprintf("%s (aktuell: %s)", identifier, contact.Name), func() {
			// Eingabefeld für neuen Namen
			entry := widget.NewEntry()
			entry.SetPlaceHolder("Neuer Name")

			form := dialog.NewForm("Kontaktnamen bearbeiten", "Speichern", "Abbrechen", []*widget.FormItem{
				widget.NewFormItem("Neuer Name", entry),
			}, func(b bool) {
				if b {
					newName := strings.TrimSpace(entry.Text)
					if newName != "" {
						contacts.UpdateContactName(identifierCopy, newName)
						dialog.ShowInformation("Erfolg", "Kontaktnamen aktualisiert.", window)
					} else {
						dialog.ShowInformation("Fehler", "Name darf nicht leer sein.", window)
					}
				}
			}, window)
			form.Show()
		})
		items = append(items, item)
	}

	// "Zurück"-Button hinzufügen
	backButton := widget.NewButton("Zurück", func() {
		showMainMenu(window)
	})

	items = append(items, backButton)
	content := container.NewVBox(items...)
	scroll := container.NewScroll(content)
	window.SetContent(scroll)
}
