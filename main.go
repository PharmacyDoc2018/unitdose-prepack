package main

import (
	"fmt"
)

func main() {
	c := initConfig()

	//c.MedProducts.RemoveProduct("NDC", "69618-0024-01")

	//err := c.MedProducts.AddProduct("diphenhydrAMINE", "25 mg", "PO CAP", "Reliable-1 Laboratories", "69618-0024-01", "00369618024014")

	//err := c.PrePackTemplates.AddTemplate("diphenhydrAMINE", "25 mg", "PO CAP", "6", 180*24*time.Hour)

	c.NonControlLog.AddEntry(0, 0, 100, "25J343", "9/30/2028")

	errSlice := c.saveData()
	if len(errSlice) > 0 {
		for _, e := range errSlice {
			fmt.Println(e.Error())
		}
	}

}
