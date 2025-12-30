package main

import (
	"fmt"
)

func main() {
	c := initConfig()

	//c.MedProducts.RemoveProduct("NDC", "69618-0024-01")

	//err := c.MedProducts.AddProduct("diphenhydrAMINE", "25 mg", "PO CAP", "Reliable-1 Laboratories", "69618-0024-01", "00369618024014")

	//err := c.PrePackTemplates.AddTemplate("diphenhydrAMINE", "25 mg", "PO CAP", "6", 180*24*time.Hour)

	c.ControlTwoLog.ControlCatagories = []string{"2"}
	c.ControlThreeToFiveLog.ControlCatagories = []string{"3", "4", "5"}
	c.NonControlLog.ControlCatagories = []string{"6"}

	fmt.Println(c.PrePackTemplates.GetMedications())

	for med := range c.MedProducts.Map {
		fmt.Println(med)

		for dose := range c.MedProducts.Map[med] {
			fmt.Println(dose)

			for form := range c.MedProducts.Map[med][dose] {
				fmt.Println(form)

				for _, prdct := range c.MedProducts.Map[med][dose][form] {
					fmt.Println(prdct.NDC)
				}
			}
		}
	}

	errSlice := c.saveData()
	if len(errSlice) > 0 {
		for _, e := range errSlice {
			fmt.Println(e.Error())
		}
	}

}
