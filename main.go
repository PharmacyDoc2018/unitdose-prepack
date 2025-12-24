package main

import "fmt"

func main() {
	c := initConfig()

	err := c.MedProducts.AddProduct("diphenhydrAMINE PO CAP", "25 mg", "PO CAP", "Reliable-1 Laboratories", "69618-0024-01", "00369618024014")
	if err != nil {
		fmt.Println(err.Error())
	}

	for key := range c.MedProducts.Map {
		fmt.Println(key)
	}

	errSlice := c.saveData()
	if len(errSlice) > 1 {
		for _, e := range errSlice {
			fmt.Println(e.Error())
		}
	}

	fmt.Println(c.medProductsPath)
	//fmt.Printf("MedProducts: %#v\n", c.MedProducts)

}
