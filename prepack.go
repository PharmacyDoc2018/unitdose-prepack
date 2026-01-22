package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/PharmacyDoc2018/unitdose-prepack/internal/barcode"
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

func (p *PrePackTemplates) GetMedications() map[string]struct{} {
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

func (p *PrePackTemplates) ListTemplates() []string {
	templateList := []string{}
	for _, template := range p.List {
		templateList = append(templateList, template.Medication)
	}

	return templateList
}

func (p *PrePackTemplates) ListNonControlTemplates() []string {
	templateList := []string{}
	for _, template := range p.List {
		if template.ControlCatagory == "6" {
			templateList = append(templateList, fmt.Sprintf("%s %s %s", template.Medication, template.Dose, template.Form))
		}
	}

	return templateList

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

type PrePackEntry struct {
	Date            time.Time
	PrePackLot      string
	PrePackTemplate PrePackTemplate
	MfgProduct      MfgProduct
	MfgLot          string
	MfgExp          string
	BarcodePath     string
	Quantity        int
}

type PrePackLog struct {
	List              []PrePackEntry
	ControlCatagories []string
	prePacktemplates  *PrePackTemplates
}

func (p *PrePackLog) AddEntry(templateIndex, productIndex, quantity int, mfgLot, mfgExp string) error {
	template := p.prePacktemplates.List[templateIndex]
	mfgProduct := p.prePacktemplates.medProducts.Map[template.Medication][template.Dose][template.Form][productIndex]

	isValidControlCat := func(tempCat string, validCats []string) bool {
		for _, cat := range validCats {
			if cat == tempCat {
				return true
			}
		}
		return false
	}(template.ControlCatagory, p.ControlCatagories)
	if !isValidControlCat {
		return fmt.Errorf("error. template %s %s %s is control catagory %s and not allowed for entry into selected log",
			template.Medication, template.Dose, template.Form, template.ControlCatagory)
	}

	if !template.Active {
		return fmt.Errorf("error. selected template is not active")
	}

	entryExp := time.Now().Add(template.BUD)

	maxSoFar := 1
	for _, entry := range p.List {
		isTodayEntry := func(e PrePackEntry) bool {
			if time.Now().Year() == e.Date.Year() &&
				time.Now().Month() == e.Date.Month() &&
				time.Now().Day() == e.Date.Day() {
				return true
			}
			return false
		}(entry)

		if isTodayEntry {
			endNum, err := strconv.Atoi(string([]rune(entry.PrePackLot)[10:]))
			if err != nil {
				return err
			}

			if endNum > maxSoFar {
				maxSoFar = endNum
			}

		}
	}

	entryLot := time.Now().Format("060102") + "-"
	if maxSoFar <= 9 {
		entryLot += "0" + strconv.Itoa(maxSoFar)
	} else {
		entryLot += strconv.Itoa(maxSoFar)
	}

	barcodePath, err := barcode.GenerateBarcode(mfgProduct.GTIN, entryExp.Format("01/02/2006"), mfgLot, entryLot)
	if err != nil {
		return err
	}

	p.List = append(p.List, PrePackEntry{
		Date:            time.Now(),
		PrePackLot:      entryLot,
		PrePackTemplate: template,
		MfgProduct:      mfgProduct,
		MfgLot:          mfgLot,
		MfgExp:          mfgExp,
		BarcodePath:     barcodePath,
		Quantity:        quantity,
	})

	return nil
}

func (p *PrePackLog) Len() int {
	return len(p.List)
}

func (c *config) LoadControlTwoLog() error {
	_, err := os.Stat(c.controlTwoPath)
	if err != nil {
		return err
	}

	controlTwoLog := PrePackLog{}
	data, err := os.ReadFile(c.controlTwoPath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &controlTwoLog)
	if err != nil {
		return err
	}

	c.ControlTwoLog = controlTwoLog

	return nil
}

func (c *config) SaveControlTwoLog() error {
	data, err := json.Marshal(c.ControlTwoLog)
	if err != nil {
		return err
	}

	saveFile, err := os.OpenFile(c.controlTwoPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
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

func (c *config) LoadControlThreeToFiveLog() error {
	_, err := os.Stat(c.controlThreeToFivePath)
	if err != nil {
		return err
	}

	controlThreeToFiveLog := PrePackLog{}
	data, err := os.ReadFile(c.controlThreeToFivePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &controlThreeToFiveLog)
	if err != nil {
		return err
	}

	c.ControlThreeToFiveLog = controlThreeToFiveLog

	return nil
}

func (c *config) SaveControlThreeToFiveLog() error {
	data, err := json.Marshal(c.ControlThreeToFiveLog)
	if err != nil {
		return err
	}

	saveFile, err := os.OpenFile(c.controlThreeToFivePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
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

func (c *config) LoadNonControlLog() error {
	_, err := os.Stat(c.nonControlPath)
	if err != nil {
		return err
	}

	nonControlLog := PrePackLog{}
	data, err := os.ReadFile(c.nonControlPath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &nonControlLog)
	if err != nil {
		return err
	}

	c.NonControlLog = nonControlLog

	return nil
}

func (c *config) SaveNonControlLog() error {
	data, err := json.Marshal(c.NonControlLog)
	if err != nil {
		return err
	}

	saveFile, err := os.OpenFile(c.nonControlPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
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

type PrepPersons struct {
	Map map[string]struct{}
}

type CheckPersons struct {
	Map map[string]struct{}
}
