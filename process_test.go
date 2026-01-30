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

func TestFormatMfgExpDate(t *testing.T) {
	tests := []struct {
		name           string
		date           string
		expectedResult string
		expectedErr    bool
	}{
		{
			"remove leading zeros",
			"09/04/2025",
			"9/4/2025",
			false,
		},
		{
			"extend short date",
			"3/18/22",
			"3/18/2022",
			false,
		},
		{
			"unchanged",
			"5/25/2021",
			"5/25/2021",
			false,
		},
		{
			"giberish",
			"aksdfajeif",
			"",
			true,
		},
		{
			"empty",
			"",
			"",
			true,
		},
		{
			"replace hyphens",
			"4-26-2015",
			"4/26/2015",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := formatMfgExpDate(tt.date)
			if result != tt.expectedResult {
				t.Errorf("fail. Expected: %s. Actual: %s", tt.expectedResult, result)
			}

			if tt.expectedErr && err == nil {
				t.Errorf("fail. Expected error. Actual: nil")
			}
		})
	}
}
