package main

type MfgProduct struct {
	MfgName string `json:"mfg_name"`
	NDC     string `json:"ndc"`
	GTIN    string `json:"gtin"`
}

type MedProducts struct {
	Map map[string]map[string]map[string]map[MfgProduct]struct{} `json:"map"`
}

func (m *MedProducts) AddProduct(medication, dose, form, mfgName, NDC, GTIN string) error {
	mfgProduct := MfgProduct{
		MfgName: mfgName,
		NDC:     NDC,
		GTIN:    GTIN,
	}

	if m.Map[medication] == nil {
		m.Map[medication] = map[string]map[string]map[MfgProduct]struct{}{}
	}

	if m.Map[medication][dose] == nil {
		m.Map[medication][dose] = map[string]map[MfgProduct]struct{}{}
	}

	if m.Map[medication][dose][form] == nil {
		m.Map[medication][dose][form] = map[MfgProduct]struct{}{}
	}

	m.Map[medication][dose][form][mfgProduct] = struct{}{}

	return nil
}
