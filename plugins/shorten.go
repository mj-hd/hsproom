package plugins

import "regexp"

func shorten(params []interface{}) interface{} {

	tar, ok := params[0].(string)

	if !ok {
		return "See other..."
	}

	re := regexp.MustCompile("(?m)[\\s]+")

	return re.ReplaceAllString(tar, "  ")
}
