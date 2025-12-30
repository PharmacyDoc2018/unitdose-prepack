package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type PrePackTemplate struct {
	Medication      string
	Dose            string
	Form            string
	BUD             time.Duration
	ControlCatagory string
	Active          bool
}

type PrePackTemplates struct {
	List        []PrePackTemplate
	medProducts *MedProducts
}

func (p *PrePackTemplates) GetMfgProducts(i int) []MfgProduct {
	return p.medProducts.Map[p.List[i].Medication][p.List[i].Dose][p.List[i].Form]
}

func (p *PrePackTemplates) getMedications() map[string]struct{} {
	returnMap := map[string]struct{}{}

	for key := range p.medProducts.Map {
		returnMap[key] = struct{}{}
	}

	return returnMap
}

func (p *PrePackTemplates) AddTemplate(medication, dose, form, controlCatagory string, BUD time.Duration) error {
	validControlCatagories := map[string]struct{}{
		"1": struct{}{},
		"2": struct{}{},
		"3": struct{}{},
		"4": struct{}{},
		"5": struct{}{},
		"6": struct{}{},
	}

	if _, ok := p.medProducts.Map[medication][dose][form]; !ok {
		return fmt.Errorf("error. no mfg product found for %s %s %s", medication, dose, form)
	}

	if _, ok := validControlCatagories[controlCatagory]; !ok {
		return fmt.Errorf("error. %s is not a valid control catagory", controlCatagory)
	}

	p.List = append(p.List, PrePackTemplate{
		Medication:      medication,
		Dose:            dose,
		Form:            form,
		BUD:             BUD,
		ControlCatagory: controlCatagory,
		Active:          true,
	})

	return nil
}

func (c *config) LoadPrePackTemplates() error {
	_, err := os.Stat(c.prePackTemplatesPath)
	if err != nil {
		return err
	}

	prePackTemplates := PrePackTemplates{}
	data, err := os.ReadFile(c.prePackTemplatesPath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &prePackTemplates)
	if err != nil {
		return err
	}

	c.PrePackTemplates = prePackTemplates

	return nil

}

func (c *config) SavePrePackTemplates() error {
	data, err := json.Marshal(c.PrePackTemplates)
	if err != nil {
		return err
	}

	saveFile, err := os.OpenFile(c.prePackTemplatesPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer saveFile.Close()

	_, err = saveFile.Write(data)
	if err != nil {
		return err
	}

	return nil
}

//func (p *PrePackTemplates) AddTemplate()

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
