package plugins

import (
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

func replaceSourceTag(params []interface{}) interface{} {

	input, ok := params[0].(string)
	if !ok {
		return "See Other..."
	}

	source, ok := params[1].(string)
	if !ok {
		return "See Other..."
	}

	if source == "" {
		return strings.Replace(input, "[sourcecode]", "", 1)
	} else {
		return strings.Replace(input, "[sourcecode]", "<pre id='sourcecode'>\n"+bluemonday.StrictPolicy().Sanitize(source)+"\n</pre>", 1)
	}
}
