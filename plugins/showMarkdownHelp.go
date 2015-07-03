package plugins

import (
	"html/template"
	"io/ioutil"
	"os"

	"../config"
)

func showMarkdownHelp(params []interface{}) interface{} {

	fp, err := os.Open(config.TemplatesPath + "plugins/markdownHelp.tmpl")
	if err != nil {
		return "See Other..."
	}

	byt, err := ioutil.ReadAll(fp)
	if err != nil {
		return "See Other..."
	}

	return template.HTML(string(byt))
}
