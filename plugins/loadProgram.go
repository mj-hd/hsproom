package plugins

import (
	"io/ioutil"
	"os"

	"hsproom/config"
	"html/template"
)

func loadProgram() interface{} {

	fp, err := os.Open(config.TemplatesPath + "plugins/loadProgram.tmpl")
	if err != nil {
		return ""
	}

	byt, err := ioutil.ReadAll(fp)
	if err != nil {
		return ""
	}

	return template.JS(string(byt))
}
