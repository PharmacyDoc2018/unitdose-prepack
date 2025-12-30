package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	medProductsPath       string
	MedProducts           MedProducts
	prePackTemplatesPath  string
	PrePackTemplates      PrePackTemplates
	ControlTwoLog         PrePackLog
	ControlThreeToFiveLog PrePackLog
	NonControlLog         PrePackLog
}

func initConfig() *config {
	c := &config{}

	godotenv.Load(".env")
	c.medProductsPath = os.Getenv("MED_PRODUCTS_PATH")
	c.prePackTemplatesPath = os.Getenv("PREPACK_TEMPLATES_PATH")

	c.MedProducts.Map = map[string]map[string]map[string][]MfgProduct{}
	c.PrePackTemplates.List = []PrePackTemplate{}

	errorSlice := c.loadData()
	if len(errorSlice) > 0 {
		for _, e := range errorSlice {
			fmt.Println(e.Error())
		}
	}

	c.PrePackTemplates.medProducts = &c.MedProducts

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

	return errorSlice
}
