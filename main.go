package main

import "fmt"

func main() {
	c := initConfig()

	//c.MedProducts.RemoveProduct("NDC", "69618-0024-01")

	//err := c.MedProducts.AddProduct("diphenhydrAMINE", "25 mg", "PO CAP", "Reliable-1 Laboratories", "69618-0024-01", "00369618024014")

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
