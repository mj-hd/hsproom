package plugins

import (
	"io/ioutil"
	"os"

	"hsproom/config"
	"html/template"
)

func loadProgramHGIMG4() interface{} {

	fp, err := os.Open(config.TemplatesPath + "plugins/loadProgramHGIMG4.tmpl")
	if err != nil {
		return ""
	}

	byt, err := ioutil.ReadAll(fp)
	if err != nil {
		return ""
	}

	return template.JS(string(byt))
}

func loadProgramHSP3Dish() interface{} {

	fp, err := os.Open(config.TemplatesPath + "plugins/loadProgramHSP3Dish.tmpl")
	if err != nil {
		return ""
	}

	byt, err := ioutil.ReadAll(fp)
	if err != nil {
		return ""
	}

	return template.JS(string(byt))
}
