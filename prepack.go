package main

import "time"

type PrePackTemplate struct {
	Medication      string
	Dose            string
	Form            string
	Product         MfgProduct
	BUD             time.Duration
	ControlCatagory string
}

type PrePackTemplates struct {
	Map map[PrePackTemplate]struct{}
}

type PrePackEntry struct {
	Date        time.Time
	Lot         string
	Medication  PrePackTemplate
	barcodePath string
	Quantity    int
}

type PrePackLog struct {
	List []PrePackEntry
}

type PrepPersons struct {
	Map map[string]struct{}
}

type CheckPersons struct {
	Map map[string]struct{}
}
