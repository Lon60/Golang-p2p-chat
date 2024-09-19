package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func StartUI() {
	myApp := app.New()

	myWindow := myApp.NewWindow("Peer-to-Peer Chat")

	showMainMenu(myWindow)

	myWindow.Resize(fyne.NewSize(400, 400))
	myWindow.ShowAndRun()
}

func showMainMenu(window fyne.Window) {
	label := widget.NewLabel("Welcome to Peer-to-Peer Chat")
	content := container.NewVBox(
		label,
		widget.NewButton("Send Contact Request", func() {
			sendContactRequestUI(window)
		}),
		widget.NewButton("View Received Contact Requests", func() {
			viewReceivedRequestsUI(window)
		}),
		widget.NewButton("View Chats", func() {
			viewChatsUI(window)
		}),
		widget.NewButton("Edit Contact Name", func() {
			editContactNameUI(window)
		}),
		widget.NewButton("Exit", func() {
			window.Close()
		}),
	)

	window.SetContent(content)
}
