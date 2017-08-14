package plugins

import (
	"io/ioutil"
	"os"

	"html/template"

	"../config"
)

func loadProgram(params []interface{}) interface{} {

	runtime, ok := params[0].(string)

	if !ok {
		return "See Other..."
	}

	switch runtime {
	case "HSP3Dish", "HGIMG4":
	default:
		return "See Other..."
	}

	version, ok := params[1].(string)
	if !ok {
		return "See Other..."
	}

	if !config.IsValidRuntimeVersion(version) {
		return "See Other..."
	}

	fp, err := os.Open(config.TemplatesPath + "plugins/loadProgram" + runtime + "/" + version + ".tmpl")
	if err != nil {
		return ""
	}

	byt, err := ioutil.ReadAll(fp)
	if err != nil {
		return ""
	}

	return template.JS(string(byt))
}
