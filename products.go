package main

type MfgProduct struct {
	MfgName string
	NDC     string
	GTIN    string
}

type MedProducts struct {
	Map map[string]map[string]map[string]map[MfgProduct]struct{}
	//medication -> dose -> form -> MfgProduct
}

func (m *MedProducts) AddProduct(medication, dose, form, mfgName, NDC, GTIN string) error {
	mfgProduct := MfgProduct{
		MfgName: mfgName,
		NDC:     NDC,
		GTIN:    GTIN,
	}

	m.Map[medication][dose][form][mfgProduct] = struct{}{}
	return nil
}
