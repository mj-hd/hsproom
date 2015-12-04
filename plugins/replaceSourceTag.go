package plugins

import (
	"html/template"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

func replaceSourceTag(params []interface{}) interface{} {

	input, ok := params[0].(string)
	if !ok {
		input_html, ok := params[0].(template.HTML)
		if !ok {
			return "See Other... P0"
		}

		input = string(input_html)
	}

	source, ok := params[1].(string)
	if !ok {
		return "See Other... P1"
	}

	if source == "" {
		return strings.Replace(input, "[sourcecode]", "", 1)
	} else {
		return strings.Replace(input, "[sourcecode]", "<pre id='sourcecode' class='brush: hsp;'>\n"+bluemonday.StrictPolicy().Sanitize(source)+"\n</pre>", 1)
	}
}
