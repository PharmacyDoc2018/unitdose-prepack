package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
)

func (c *config) initHomeScreen() fyne.Window {
	w := c.App.NewWindow("Unitdose Prepack")
	w.Resize(fyne.Size{
		Width:  500,
		Height: 500,
	})

	medicationEntry := widget.NewSelectEntry(c.PrePackTemplates.ListTemplates())

	windowLabel := widget.NewLabel("Unitdose Prepack")

	form := widget.NewForm(
		widget.NewFormItem("Medication", medicationEntry),
		widget.NewFormItem("Dose", widget.NewEntry()),
	)
	w.SetContent(container.NewVBox(
		windowLabel,
		form,
	))

	return w
}

func (c *config) addNewEntryWindow() fyne.Window {
	// Incomplete - current code in for reference only
	w := c.App.NewWindow("Add Entry")
	w.Resize(fyne.Size{
		Width:  500,
		Height: 500,
	})

	medicationEntry := widget.NewSelectEntry(c.PrePackTemplates.ListTemplates())

	hello := widget.NewLabel("Hello Fyne!")
	form := widget.NewForm(
		widget.NewFormItem("Medication", medicationEntry),
		widget.NewFormItem("Dose", widget.NewEntry()),
	)
	w.SetContent(container.NewVBox(
		hello,
		widget.NewButton("Hi!", func() {
			hello.SetText("Welcome :)")
		}),
		form,
	))

	return w

}
