package main

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/joho/godotenv"
)

type config struct {
	medProductsPath        string
	MedProducts            MedProducts
	prePackTemplatesPath   string
	PrePackTemplates       PrePackTemplates
	controlTwoPath         string
	ControlTwoLog          PrePackLog
	controlThreeToFivePath string
	ControlThreeToFiveLog  PrePackLog
	nonControlPath         string
	NonControlLog          PrePackLog
	App                    fyne.App
}

func initConfig() *config {
	c := &config{}

	godotenv.Load(".env")
	c.medProductsPath = os.Getenv("MED_PRODUCTS_PATH")
	c.prePackTemplatesPath = os.Getenv("PREPACK_TEMPLATES_PATH")

	c.controlTwoPath = os.Getenv("C_2_PATH")
	c.controlThreeToFivePath = os.Getenv("C_3_TO_5_PATH")
	c.nonControlPath = os.Getenv("NON_CONTROL_PATH")

	c.MedProducts.Map = map[string]map[string]map[string][]MfgProduct{}
	c.PrePackTemplates.List = []PrePackTemplate{}

	errorSlice := c.loadData()
	if len(errorSlice) > 0 {
		for _, e := range errorSlice {
			fmt.Println(e.Error())
		}
	}

	c.PrePackTemplates.medProducts = &c.MedProducts

	c.ControlTwoLog.prePacktemplates = &c.PrePackTemplates
	c.ControlThreeToFiveLog.prePacktemplates = &c.PrePackTemplates
	c.NonControlLog.prePacktemplates = &c.PrePackTemplates

	c.App = app.New()

	return c
}

func (c *config) saveData() []error {

	errorSlice := []error{}

	err := c.SaveMedProducts()
	if err != nil {
		errorSlice = append(errorSlice, err)
	}

	err = c.SavePrePackTemplates()
	if err != nil {
		errorSlice = append(errorSlice, err)
	}

	err = c.SaveControlTwoLog()
	if err != nil {
		errorSlice = append(errorSlice, err)
	}

	err = c.SaveControlThreeToFiveLog()
	if err != nil {
		errorSlice = append(errorSlice, err)
	}

	err = c.SaveNonControlLog()
	if err != nil {
		errorSlice = append(errorSlice, err)
	}

	return errorSlice

}

func (c *config) loadData() []error {
	errorSlice := []error{}

	err := c.LoadMedProducts()
	if err != nil {
		errorSlice = append(errorSlice, err)
	}

	err = c.LoadPrePackTemplates()
	if err != nil {
		errorSlice = append(errorSlice, err)
	}

	err = c.LoadControlTwoLog()
	if err != nil {
		errorSlice = append(errorSlice, err)
	}

	err = c.LoadControlThreeToFiveLog()
	if err != nil {
		errorSlice = append(errorSlice, err)
	}

	err = c.LoadNonControlLog()
	if err != nil {
		errorSlice = append(errorSlice, err)
	}

	return errorSlice
}
