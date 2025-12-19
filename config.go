package main

type config struct {
	MedProducts           MedProducts
	PrePackTemplates      PrePackTemplates
	ControlTwoLog         PrePackLog
	ControlThreeToFiveLog PrePackLog
	NonControlLog         PrePackLog
}

func initConfig() *config {
	c := config{}

	c.MedProducts.Map = map[string]map[string]map[string]map[MfgProduct]struct{}{}

	c.PrePackTemplates.Map = map[PrePackTemplate]struct{}{}

	return &c
}
