package ui

import (
	"Golang-p2p-chat/client"
	"Golang-p2p-chat/config"
	"Golang-p2p-chat/models"
	"Golang-p2p-chat/server"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"strings"
)

func sendContactRequestUI(window fyne.Window) {
	ipEntry := widget.NewEntry()
	ipEntry.SetPlaceHolder("IP-Adresse des Kontakts")
	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder("Port des Kontakts")

	formItems := []*widget.FormItem{
		widget.NewFormItem("IP-Adresse", ipEntry),
		widget.NewFormItem("Port", portEntry),
	}

	form := dialog.NewForm("Kontaktanfrage senden", "Senden", "Abbrechen", formItems, func(b bool) {
		if b {
			ip := strings.TrimSpace(ipEntry.Text)
			port := strings.TrimSpace(portEntry.Text)

			contactRequest := models.ContactRequest{
				Name: config.LocalUserName,
				IP:   getLocalIP(),
				Port: server.PORT,
			}
			err := client.SendContactRequest(ip, port, contactRequest)
			if err != nil {
				dialog.ShowError(err, window)
			} else {
				dialog.ShowInformation("Erfolg", "Kontaktanfrage gesendet.", window)
			}
		}
	}, window)

	form.Show()
}
