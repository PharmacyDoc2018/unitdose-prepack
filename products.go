package main

import "time"

type MfgProduct struct {
	MfgName string
	NDC     string
	GTIN    string
	Exp     time.Time
	Lot     string
}

type MedProducts struct {
	Map map[string]map[string]map[string]map[MfgProduct]struct{}
	//medication -> dose -> form -> MfgProduct
}
