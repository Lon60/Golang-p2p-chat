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
	label := widget.NewLabel("Willkommen zum Peer-to-Peer Chat")
	content := container.NewVBox(
		label,
		widget.NewButton("Kontaktanfrage senden", func() {
			sendContactRequestUI(window)
		}),
		widget.NewButton("Empfangene Kontaktanfragen anzeigen", func() {
			viewReceivedRequestsUI(window)
		}),
		widget.NewButton("Chats anzeigen", func() {
			viewChatsUI(window)
		}),
		widget.NewButton("Kontaktnamen bearbeiten", func() {
			editContactNameUI(window)
		}),
		widget.NewButton("Beenden", func() {
			window.Close()
		}),
	)

	window.SetContent(content)
}
