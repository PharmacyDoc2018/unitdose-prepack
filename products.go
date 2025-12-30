package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type MfgProduct struct {
	MfgName string `json:"mfg_name"`
	NDC     string `json:"ndc"`
	GTIN    string `json:"gtin"`
}

type MedProducts struct {
	Map map[string]map[string]map[string][]MfgProduct `json:"map"`
}

func (m *MedProducts) AddProduct(medication, dose, form, mfgName, NDC, GTIN string) error {
	mfgProduct := MfgProduct{
		MfgName: mfgName,
		NDC:     NDC,
		GTIN:    GTIN,
	}

	if m.Map[medication] == nil {
		m.Map[medication] = map[string]map[string][]MfgProduct{}
	}

	if m.Map[medication][dose] == nil {
		m.Map[medication][dose] = map[string][]MfgProduct{}
	}

	if m.Map[medication][dose][form] == nil {
		m.Map[medication][dose][form] = []MfgProduct{}
	}

	m.Map[medication][dose][form] = append(m.Map[medication][dose][form], mfgProduct)

	return nil
}

func (m *MedProducts) RemoveProduct(idType, id string) error {
	markedMed := ""
	markedDose := ""
	markedForm := ""
	markedIndex := 0

	switch idType {
	case "NDC", "ndc":
		for med := range m.Map {
			for dose := range m.Map[med] {
				for form := range m.Map[med][dose] {
					for i, prdct := range m.Map[med][dose][form] {
						if prdct.NDC == id {
							markedMed = med
							markedDose = dose
							markedForm = form
							markedIndex = i
							break
						}
					}
					if markedMed != "" {
						break
					}
				}
				if markedMed != "" {
					break
				}
			}
			if markedMed != "" {
				break
			}
		}

		if markedMed == "" {
			return fmt.Errorf("error. product not found with NDC: %s", id)
		}

	case "GTIN", "gtin":
		for med := range m.Map {
			for dose := range m.Map[med] {
				for form := range m.Map[med][dose] {
					for i, prdct := range m.Map[med][dose][form] {
						if prdct.GTIN == id {
							markedMed = med
							markedDose = dose
							markedForm = form
							markedIndex = i
							break
						}
					}
					if markedMed != "" {
						break
					}
				}
				if markedMed != "" {
					break
				}
			}
			if markedMed != "" {
				break
			}
		}

		if markedMed == "" {
			return fmt.Errorf("error. product not found with GTIN: %s", id)
		}

	default:
		return fmt.Errorf("error. Invalid id type")
	}

	m.Map[markedMed][markedDose][markedForm] = append(m.Map[markedMed][markedDose][markedForm][:markedIndex], m.Map[markedMed][markedDose][markedForm][markedIndex+1:]...)

	if len(m.Map[markedMed][markedDose][markedForm]) == 0 {
		delete(m.Map[markedMed][markedDose], markedForm)
	}

	if len(m.Map[markedMed][markedDose]) == 0 {
		delete(m.Map[markedMed], markedDose)
	}

	if len(m.Map[markedMed]) == 0 {
		delete(m.Map, markedMed)
	}

	return nil
}

func (c *config) LoadMedProducts() error {
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

	if medProducts.Map == nil {
		c.MedProducts.Map = map[string]map[string]map[string][]MfgProduct{}
	} else {
		c.MedProducts = medProducts
	}

	return nil
}

func (c *config) SaveMedProducts() error {
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
}
