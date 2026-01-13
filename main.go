package main

import (
	"fmt"
)

func main() {
	c := initConfig()
	w := c.initHomeScreen()

	w.ShowAndRun()

	errSlice := c.saveData()
	if len(errSlice) > 0 {
		for _, e := range errSlice {
			fmt.Println(e.Error())
		}
	}

}
