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

	n, d, f, err := templates.ValidateNDC(testTemplateName, testNDC)
	if err != nil {
		t.Errorf("NDC validation failed: %s", err.Error())
	}

	if n != "diphenhydrAMINE" {
		t.Errorf("medication name does not match: %s", n)
	}

	if d != "25 mg" {
		t.Errorf("dose does not match %s", d)
	}

	if f != "PO CAP" {
		t.Errorf("form does not match: %s", f)
	}
}
