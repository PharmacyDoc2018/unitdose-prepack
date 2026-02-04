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

type WindowType int

const (
	Main WindowType = iota
	AddEntry
)

type guiState struct {
	nonControlLogEntriesBinding binding.List[any]
	addEntryButton              *widget.Button
	editEntryButton             *widget.Button
	removeEntryButton           *widget.Button
	noncontrolLogSelectedID     widget.ListItemID
	cThreeToFiveLogSelectedID   widget.ListItemID
	cTwoLogSelectedID           widget.ListItemID
	currentTab                  int
	currentWindowType           WindowType
	tabs                        *container.AppTabs
	addEntryForm                *widget.Form
	nonControlLogTable          *fyne.Container
}

type myTheme struct{}

func (myTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	if n == theme.ColorNameBackground {
		return color.NRGBA{R: 240, G: 240, B: 240, A: 255}
	}

	if n == theme.ColorNameButton {
		return color.NRGBA{R: 173, G: 216, B: 230, A: 255}
	}

	if n == theme.ColorNameDisabledButton {
		return color.NRGBA{R: 210, G: 210, B: 210, A: 255}
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

func initGUIState() *guiState {
	s := guiState{}
	s.noncontrolLogSelectedID = -1
	s.cThreeToFiveLogSelectedID = -1
	s.cTwoLogSelectedID = -1
	s.currentTab = 0
	s.currentWindowType = Main

	return &s
}

func (c *config) initHomeScreen() fyne.Window {
	s := initGUIState()
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

	addEntryButton := widget.NewButton("Add Entry", func() {
		newW := addNewEntryWindow(c, s)
		s.currentWindowType = AddEntry
		s.stateChange()
		newW.Show()
	})
	s.addEntryButton = addEntryButton

	editEntryButton := widget.NewButton("Edit Entry", func() {})
	editEntryButton.Disable()
	s.editEntryButton = editEntryButton

	removeEntryButton := widget.NewButton("Remove Entry", func() {
		switch s.currentTab {
		case 0: //-- NonControled Tab
			val, _ := s.nonControlLogEntriesBinding.GetValue(s.noncontrolLogSelectedID)
			entry := val.(PrePackEntry)
			c.NonControlLog.RemoveEntry(entry.PrePackLot)
			err := s.nonControlLogEntriesBinding.Remove(entry)
			if err != nil {
				fmt.Println(err.Error())
			}
			s.stateChange()
		}
	})
	removeEntryButton.Disable()
	s.removeEntryButton = removeEntryButton

	bottomButtonsContainer := container.NewHBox(addEntryButton, s.editEntryButton, s.removeEntryButton)

	s.nonControlLogTable = createNonControlLogTable(c, s)
	nonControlLogScroll := container.NewScroll(s.nonControlLogTable)
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
	tabs.SelectedIndex()
	tabs.OnSelected = func(currentTab *container.TabItem) {
		s.currentTab = tabs.SelectedIndex()
		s.stateChange()
	}
	s.tabs = tabs

	w.SetContent(container.NewVBox(
		windowLabel,
		tabs,
		bottomButtonsContainer,
	))

	return w
}

func addNewEntryWindow(c *config, s *guiState) fyne.Window {
	w := c.App.NewWindow("Add Entry")
	w.Resize(fyne.Size{
		Width:  500,
		Height: 500,
	})

	w.SetOnClosed(func() {
		s.currentWindowType = Main
		s.stateChange()
	})

	medicationEntryOptions := []string{}
	switch s.currentTab {
	case 0: //-- NonControl Tab
		medicationEntryOptions = c.PrePackTemplates.ListNonControlTemplates()

	}
	medicationEntry := widget.NewSelectEntry(medicationEntryOptions)
	ndcEntry := widget.NewEntry()
	mfgLotEntry := widget.NewEntry()
	mfgExpEntry := widget.NewEntry()
	quantityEntry := widget.NewEntry()

	form := widget.NewForm(
		widget.NewFormItem("Template", medicationEntry),
		widget.NewFormItem("NDC", ndcEntry),
		widget.NewFormItem("Mfg Lot", mfgLotEntry),
		widget.NewFormItem("Mfg Exp", mfgExpEntry),
		widget.NewFormItem("Quantity", quantityEntry),
	)
	form.SubmitText = "Add Entry"
	form.OnSubmit = func() {
		templateIndex, productIndex, err := c.PrePackTemplates.ValidateNDC(medicationEntry.Text, ndcEntry.Text)
		if err != nil {
			fmt.Println(err.Error())
		}

		qty, err := strconv.Atoi(quantityEntry.Text)
		if err != nil {
			fmt.Println(err.Error())
			w.Close()
			return
		}
		newEntry := PrePackEntry{}
		switch s.currentTab {
		case 0: //-- NonControl Tab
			newEntry, err = c.NonControlLog.AddEntry(templateIndex, productIndex, qty, mfgLotEntry.Text, mfgExpEntry.Text)
			if err != nil {
				fmt.Println(err.Error())
				w.Close()
				return
			}
		}
		s.nonControlLogEntriesBinding.Append(newEntry)

		w.Close()
	}
	form.OnCancel = func() {
		w.Close()
	}
	s.addEntryForm = form

	w.SetContent(container.NewVBox(
		form,
	))

	return w

}

func createNonControlLogTable(c *config, s *guiState) *fyne.Container {
	nonControlLogEntriesBinding := binding.NewUntypedList()
	for _, e := range c.NonControlLog.List {
		nonControlLogEntriesBinding.Append(e)
	}
	s.nonControlLogEntriesBinding = nonControlLogEntriesBinding
	nonControlLogObj := widget.NewListWithData(
		s.nonControlLogEntriesBinding,
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

	nonControlLogObj.OnSelected = func(id widget.ListItemID) {
		s.noncontrolLogSelectedID = id
		s.stateChange()
	}

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

func updateAddEntryButtonStats(s *guiState) {
	if s.currentWindowType != Main {
		s.addEntryButton.Disable()
	} else {
		s.addEntryButton.Enable()
	}
}

func updateEditEntryButtonState(s *guiState) {
	if s.currentWindowType != Main {
		s.editEntryButton.Disable()
		return
	}
	switch s.currentTab {
	case 0: //-- NonControled Tab
		if s.noncontrolLogSelectedID == -1 {
			s.editEntryButton.Disable()
		} else {
			s.editEntryButton.Enable()
		}

	case 1: //-- CIII - V Tab
		if s.cThreeToFiveLogSelectedID == -1 {
			s.editEntryButton.Disable()
		} else {
			s.editEntryButton.Enable()
		}

	case 2: //-- CII Tab
		if s.cTwoLogSelectedID == -1 {
			s.editEntryButton.Disable()
		} else {
			s.editEntryButton.Enable()
		}
	}
}

func updateRemoveEntryButtonState(s *guiState) {
	if s.currentWindowType != Main {
		s.editEntryButton.Disable()
		return
	}
	switch s.currentTab {
	case 0: //-- NonControled Tab
		if s.noncontrolLogSelectedID == -1 {
			s.removeEntryButton.Disable()
		} else {
			s.removeEntryButton.Enable()
		}

	case 1: //-- CIII - V Tab
		if s.cThreeToFiveLogSelectedID == -1 {
			s.removeEntryButton.Disable()
		} else {
			s.removeEntryButton.Enable()
		}

	case 2: //-- CII Tab
		if s.cTwoLogSelectedID == -1 {
			s.removeEntryButton.Disable()
		} else {
			s.removeEntryButton.Enable()
		}
	}
}

func (s *guiState) stateChange() {
	updateAddEntryButtonStats(s)
	updateEditEntryButtonState(s)
	updateRemoveEntryButtonState(s)

	s.nonControlLogTable.Refresh()
}
