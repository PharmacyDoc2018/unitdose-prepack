package main

import (
	"fmt"
	"strings"
	"time"
)

func formatNDC(ndc string) (string, error) {
	if ndc == "" {
		return "", fmt.Errorf("error. NDC is blank")
	}

	ndcSplit := strings.Split(ndc, "-")

	for _, slice := range ndcSplit {
		if !isAllNumeric(slice) {
			return "", fmt.Errorf("error. NDC cannot contain letters")
		}
	}

	switch len(ndcSplit) {
	case 3: //-- NDC contained 2 hyphens
		if len(ndcSplit[0]) > 5 {
			return "", fmt.Errorf("error. Invalid NDC format")
		}
		for len(ndcSplit[0]) < 5 {
			ndcSplit[0] = "0" + ndcSplit[0]
		}

		if len(ndcSplit[1]) > 4 {
			return "", fmt.Errorf("error. Invalid NDC format")
		}
		for len(ndcSplit[1]) < 4 {
			ndcSplit[1] = "0" + ndcSplit[1]
		}

		if len(ndcSplit[2]) > 2 {
			return "", fmt.Errorf("error. Invalid NDC format")
		}
		for len(ndcSplit[2]) < 2 {
			ndcSplit[2] = "0" + ndcSplit[2]
		}

	case 2: //-- NDC only contained 1 hyphen
		tempJoinNDC := strings.Join(ndcSplit, "")
		if len(tempJoinNDC) == 11 {
			ndcSplit[0] = tempJoinNDC[:5]
			ndcSplit[1] = tempJoinNDC[5:9]
			ndcSplit[2] = tempJoinNDC[9:]
		} else {
			return "", fmt.Errorf("error. Invalid NDC format")
		}

	case 1: //-- no hypens in NDC
		if len(ndcSplit[0]) == 11 {
			tempNDC := ndcSplit[0]
			ndcSplit[0] = tempNDC[:5]
			ndcSplit[1] = tempNDC[5:9]
			ndcSplit[2] = tempNDC[9:]
		} else {
			return "", fmt.Errorf("error. Invalid NDC format")
		}

	default:
		return "", fmt.Errorf("error. Invalid NDC format")
	}

	formattedNDC := strings.Join(ndcSplit, "-")
	return formattedNDC, nil
}

func isAllNumeric(inpt string) bool {
	if inpt == "" {
		return false
	}

	for _, r := range inpt {
		if r < '0' || r > '9' {
			return false
		}
	}

	return true
}

func formatMfgExpDate(date string) (string, error) {
	dateFormats := []string{
		"1/2/2006",
		"01/02/2006",
		"1/2/06",
		"01/02/06",
		"1-2-2006",
		"01-02-2006",
		"1-2/06",
		"01-02-06",
	}

	var dateTime time.Time
	var err error
	for _, format := range dateFormats {
		dateTime, err = time.Parse(format, date)
		if err != nil {
			continue
		} else {
			break
		}
	}

	if dateTime.IsZero() {
		return "", fmt.Errorf("error. unable to %s as a date", date)
	}

	return dateTime.Format("1/2/2006"), nil
}
