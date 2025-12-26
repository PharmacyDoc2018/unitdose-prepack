package main

import "fmt"

func main() {
	c := initConfig()

	for med := range c.MedProducts.Map {
		fmt.Println(med)

		for dose := range c.MedProducts.Map[med] {
			fmt.Println(dose)
		}
	}

	errSlice := c.saveData()
	if len(errSlice) > 0 {
		for _, e := range errSlice {
			fmt.Println(e.Error())
		}
	}

}
