package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	medProductsPath       string
	MedProducts           MedProducts
	PrePackTemplates      PrePackTemplates
	ControlTwoLog         PrePackLog
	ControlThreeToFiveLog PrePackLog
	NonControlLog         PrePackLog
}

func initConfig() *config {
	c := &config{}

	godotenv.Load(".env")
	c.medProductsPath = os.Getenv("MED_PRODUCTS_PATH")

	c.MedProducts.Map = map[string]map[string]map[string]map[MfgProduct]struct{}{}
	c.PrePackTemplates.Map = map[PrePackTemplate]struct{}{}

	errorSlice := c.loadData()
	if len(errorSlice) > 0 {
		for _, e := range errorSlice {
			fmt.Println(e.Error())
		}
	}

	return c
}

func (c *config) saveData() []error {

	errorSlice := []error{}

	//-- Save c.MedProducts
	err := func(c *config) error {
		data, err := json.Marshal(c.MedProducts)
		if err != nil {
			return err
		}

		saveFile, err := os.OpenFile(c.medProductsPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer saveFile.Close()

		_, err = saveFile.Write(data)
		if err != nil {
			return err
		}

		return nil

	}(c)
	if err != nil {
		errorSlice = append(errorSlice, err)
	}

	return errorSlice

}

func (c *config) loadData() []error {
	errorSlice := []error{}

	//-- Load MedProducts
	err := func(c *config) error {
		_, err := os.Stat(c.medProductsPath)
		if err != nil {
			return err
		}

		medProducts := MedProducts{}
		data, err := os.ReadFile(c.medProductsPath)
		if err != nil {
			return err
		}

		err = json.Unmarshal(data, &medProducts)
		if err != nil {
			return err
		}

		c.MedProducts = medProducts
		return nil
	}(c)
	if err != nil {
		errorSlice = append(errorSlice, err)
	}

	return errorSlice
}
