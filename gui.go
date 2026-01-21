package main

import (
	"fmt"
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type myTheme struct{}

func (myTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	if n == theme.ColorNameBackground {
		return color.NRGBA{R: 245, G: 245, B: 245, A: 255}
	}
	return theme.DefaultTheme().Color(n, v)
}

func (myTheme) Font(s fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(s)
}
func (myTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}
func (myTheme) Size(n fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(n)
}

func (c *config) initHomeScreen() fyne.Window {
	c.App.Settings().SetTheme(&myTheme{})

	w := c.App.NewWindow("Unitdose Prepack")
	w.Resize(fyne.Size{
		Width:  1000,
		Height: 1000,
	})

	windowLabel := canvas.NewText("Prepack Log", func(v fyne.ThemeVariant) color.Color {
		switch v {
		case 0:
			return color.Black

		case 1:
			return color.White

		default:
			return color.Black
		}

	}(c.App.Settings().ThemeVariant()))
	windowLabel.TextSize = 24

	nonControlLogTable := createNonControlLogTable(c)
	nonControlLogScroll := container.NewScroll(nonControlLogTable)
	nonControlLogScroll.SetMinSize(fyne.Size{
		Width:  800,
		Height: 400,
	})
	nonControlTab := container.NewTabItem("NonControl Log", nonControlLogScroll)

	cThreeToFiveLogObj := widget.NewLabel("CIII - V log will be displayed here")
	cThreeToFiveLogTab := container.NewTabItem("CIII - V Log", cThreeToFiveLogObj)

	cTwoLogObj := widget.NewLabel("CII log will be displayed here")
	cTwoTab := container.NewTabItem("CII Log", cTwoLogObj)

	tabs := container.NewAppTabs(nonControlTab, cThreeToFiveLogTab, cTwoTab)

	testButton := widget.NewButton("Test", func() {})

	w.SetContent(container.NewVBox(
		windowLabel,
		tabs,
		testButton,
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

func createNonControlLogTable(c *config) *fyne.Container {
	nonControlLogEntriesBinding := binding.NewUntypedList()
	for _, e := range c.NonControlLog.List {
		nonControlLogEntriesBinding.Append(e)
	}
	nonControlLogObj := widget.NewListWithData(
		nonControlLogEntriesBinding,
		func() fyne.CanvasObject {
			return newPrePackRow()
		},
		func(di binding.DataItem, obj fyne.CanvasObject) {
			val, _ := di.(binding.Untyped).Get()
			entry := val.(PrePackEntry)

			row := obj.(*fyne.Container)

			row.Objects[0].(*widget.Label).SetText(
				entry.Date.Format("2006-01-02"),
			)

			row.Objects[1].(*widget.Label).SetText(
				entry.PrePackLot,
			)

			row.Objects[2].(*widget.Label).SetText(
				fmt.Sprintf("%s %s %s", entry.PrePackTemplate.Medication, entry.PrePackTemplate.Dose, entry.PrePackTemplate.Form),
			)

			row.Objects[3].(*widget.Label).SetText(
				entry.MfgLot,
			)

			row.Objects[4].(*widget.Label).SetText(
				entry.MfgExp,
			)

			row.Objects[5].(*widget.Label).SetText(
				strconv.Itoa(entry.Quantity),
			)
		},
	)
	nonControlLogObjWithHeader := container.NewBorder(
		newPrePackHeader(),
		nil,
		nil,
		nil,
		nonControlLogObj,
	)

	return nonControlLogObjWithHeader
}

func newPrePackRow() fyne.CanvasObject {
	date := widget.NewLabel("")
	date.Alignment = fyne.TextAlignCenter

	prePackLot := widget.NewLabel("")
	prePackLot.Alignment = fyne.TextAlignCenter

	name := widget.NewLabel("")
	name.Alignment = fyne.TextAlignCenter

	mfgLot := widget.NewLabel("")
	mfgLot.Alignment = fyne.TextAlignCenter

	mfgExp := widget.NewLabel("")
	mfgExp.Alignment = fyne.TextAlignCenter

	qty := widget.NewLabel("")
	qty.Alignment = fyne.TextAlignCenter

	return container.NewGridWithColumns(
		6,
		date,
		prePackLot,
		name,
		mfgLot,
		mfgExp,
		qty,
	)
}

func newPrePackHeader() fyne.CanvasObject {
	bold := fyne.TextStyle{Bold: true}

	return container.NewGridWithColumns(
		6,
		widget.NewLabelWithStyle("Date", fyne.TextAlignCenter, bold),
		widget.NewLabelWithStyle("PrePack Lot", fyne.TextAlignCenter, bold),
		widget.NewLabelWithStyle("Name", fyne.TextAlignCenter, bold),
		widget.NewLabelWithStyle("Mfg Lot", fyne.TextAlignCenter, bold),
		widget.NewLabelWithStyle("Mfg Exp", fyne.TextAlignCenter, bold),
		widget.NewLabelWithStyle("Quantity", fyne.TextAlignCenter, bold),
	)
}
