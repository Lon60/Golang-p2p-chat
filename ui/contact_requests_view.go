package ui

import (
	"Golang-p2p-chat/client"
	"Golang-p2p-chat/contact_requests"
	"Golang-p2p-chat/contacts"
	"Golang-p2p-chat/models"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func viewReceivedRequestsUI(window fyne.Window) {
	requests := contact_requests.GetReceivedRequests()

	if len(requests) == 0 {
		dialog.ShowInformation("No Contact Requests", "No received contact requests.", window)
		return
	}

	var items []fyne.CanvasObject
	for _, req := range requests {
		reqCopy := req
		item := widget.NewButton(fmt.Sprintf("%s (%s:%s)", req.Name, req.IP, req.Port), func() {
			acceptContactRequest(reqCopy, window)
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

func acceptContactRequest(request models.ContactRequest, window fyne.Window) {
	identifier := request.IP + ":" + request.Port
	contact := models.Contact{
		Name:      request.Name,
		IP:        request.IP,
		Port:      request.Port,
		PublicKey: request.PublicKey,
		KeyObject: request.KeyObject,
	}

	contacts.AddContact(identifier, contact)

	err := client.SendContactAccepted(request)
	if err != nil {
		dialog.ShowError(err, window)
	} else {
		contact_requests.RemoveReceivedRequestByIdentifier(identifier)
		dialog.ShowInformation("Contact Accepted", "The contact has been accepted.", window)
	}
}
