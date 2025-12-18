package main

import (
	"fmt"

	"github.com/PharmacyDoc2018/unitdose-prepack/internal/barcode"
	"github.com/PharmacyDoc2018/unitdose-prepack/internal/config"
)

func main() {
	c := initConfig()

	GTIN := "00369618024014"
	prepakExp := "6/3/2026"
	prepakLot := "251217-01"
	mfgLot := "25J343"
	err := c.Barcodes.Add(GTIN, prepakExp, mfgLot, prepakLot)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func initConfig() *config.Config {
	config := config.Config{}

	config.Barcodes.Map = map[barcode.Barcode]struct{}{}

	return &config
}
