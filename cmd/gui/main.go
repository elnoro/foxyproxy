package main

import (
	"context"
	"log"

	"foxyproxy/droplets"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	fxpr "foxyproxy/cmd/fxpr/app"
)

func main() {
	fApp, err := fxpr.New()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	a := app.New()
	w := a.NewWindow("🦊Foxy Proxy🦊")

	var currentDroplet droplets.Server

	dropletStatus := widget.NewLabel("Droplet unavailable")
	dropletStartButton := widget.NewButton("Start droplet", func() {
		dropletStatus.SetText("Starting droplet...")
		s, err := fApp.StartTestServer(ctx)
		if err != nil {
			dropletStatus.SetText("Failed to start a droplet. Check settings")
			return
		}

		currentDroplet = s
		dropletStatus.SetText("Run ssh root@" + s.PublicIP)
	})

	dropletStopButton := widget.NewButton("Stop droplet", func() {
		if currentDroplet.Id == 0 {
			return
		}

		dropletStatus.SetText("Removing droplet...")
		fApp.DeleteDroplet(ctx, currentDroplet)
		currentDroplet = droplets.Server{}

		dropletStatus.SetText("Droplet unavailable")
	})

	content := container.New(layout.NewGridLayout(2), dropletStartButton, dropletStopButton, dropletStatus)
	w.SetContent(content)

	w.ShowAndRun()
}
