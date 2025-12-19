package barcode

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const urlStart = "https://barcode.tec-it.com/barcode.ashx?data="
const urlEnd = "&code=DataMatrix&translate-esc=on"

func GenerateBarcode(GTIN, prepakExp, mfgLot, prepakLot string) (string, error) {
	barcodePrefixMap := map[string]string{
		"GTIN": "01",
		"S/N":  "21",
		"Exp":  "17",
		"Lot":  "10",
	}

	prepakExp, err := parseExpDate(prepakExp)
	if err != nil {
		return "", err
	}

	barcodeData := ""
	barcodeData += barcodePrefixMap["GTIN"] + GTIN
	barcodeData += barcodePrefixMap["Exp"] + prepakExp
	barcodeData += barcodePrefixMap["Lot"] + mfgLot

	url := urlStart + barcodeData + urlEnd
	res, err := http.Get(url)
	if err != nil {
		return "", nil
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return "", fmt.Errorf("error. respose failed with status code: %d", res.StatusCode)
	}

	barcodePath := "barcodes/" + prepakLot + ".gif"
	out, err := os.Create(barcodePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		return "", err
	}

	return barcodePath, nil

}

func parseExpDate(exp string) (string, error) {
	dateFormats := []string{
		"1/2/2006",
		"01/02/2006",
		"1/2/06",
		"01/02/06",
		"1-2-2006",
		"01-02-2006",
		"1-2-06",
		"01-02-06",
	}

	expTime := time.Time{}
	var err error
	for _, dateFormat := range dateFormats {
		expTime, err = time.Parse(dateFormat, exp)
		if err == nil {
			break
		}
	}
	if err != nil {
		return "", fmt.Errorf("error. unable to parse expiration date: %s", exp)
	}

	formattedExp := expTime.Format("060102")
	return formattedExp, nil
}
