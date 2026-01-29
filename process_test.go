package main

import "testing"

func TestValidateNDC(t *testing.T) {
	products := MedProducts{}
	products.Map = map[string]map[string]map[string][]MfgProduct{}

	templates := PrePackTemplates{}
	templates.medProducts = &products

	err := products.AddProduct("diphenhydrAMINE", "25 mg", "PO CAP", "Reliable-1 Labs", "69618-0024-01", "00369618024014")
	if err != nil {
		t.Errorf("failed to add product: %s", err.Error())
	}

	testTemplateName := "diphenhydrAMINE 25 mg PO CAP"
	testNDC := "69618-0024-01"

	templateIndex, productIndex, err := templates.ValidateNDC(testTemplateName, testNDC)
	if err != nil {
		t.Errorf("NDC validation failed: %s", err.Error())
	}

	if templateIndex != 0 {
		t.Errorf("bad template index returned: %d", templateIndex)
	}

	if productIndex != 0 {
		t.Errorf("bad template index returned: %d", templateIndex)
	}
}
