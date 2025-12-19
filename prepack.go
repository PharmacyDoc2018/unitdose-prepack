package main

import "time"

type PrePackTemplate struct {
	Medication      string
	Dose            string
	Form            string
	Product         MfgProduct
	BUD             time.Duration
	ControlCatagory string
	Active          bool
}

type PrePackTemplates struct {
	Map map[PrePackTemplate]struct{}
}

func (p *PrePackTemplates) AddTemplate()

type PrePackEntry struct {
	Date        time.Time
	PrePackLot  string
	Medication  PrePackTemplate
	MfgLot      string
	MfgExp      time.Time
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
