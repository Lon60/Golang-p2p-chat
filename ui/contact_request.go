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
	ipEntry.SetPlaceHolder("Contact's IP address")
	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder("Contact's port")

	formItems := []*widget.FormItem{
		widget.NewFormItem("IP Address", ipEntry),
		widget.NewFormItem("Port", portEntry),
	}

	form := dialog.NewForm("Send Contact Request", "Send", "Cancel", formItems, func(b bool) {
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
				dialog.ShowInformation("Success", "Contact request sent.", window)
			}
		}
	}, window)

	form.Show()
}
