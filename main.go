package main

import (
	"fmt"

	"github.com/PharmacyDoc2018/unitdose-prepack/internal/barcode"
)

func main() {
	GTIN := "00369618024014"
	prepakExp := "6/3/2026"
	prepakLot := "251217-01"
	mfgLot := "25J343"
	_, err := barcode.GenerateBarcode(GTIN, prepakExp, mfgLot, prepakLot)
	if err != nil {
		fmt.Println(err.Error())
	}
}
