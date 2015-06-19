package plugins

import "strconv"

const (
	byteByte     int = 1
	byteKilobyte int = 1024 * byteByte
	byteMegabyte int = 1024 * byteKilobyte
	byteGigabyte int = 1024 * byteMegabyte
)

func byteFormat(params []interface{}) interface{} {

	byt, ok := params[0].(int)

	if !ok {
		return "See Other..."
	}

	var unit string

	switch {
	case byt >= byteGigabyte:
		unit = "GB"
		byt = byt / byteGigabyte
	case byt >= byteMegabyte:
		unit = "MB"
		byt = byt / byteMegabyte
	case byt >= byteKilobyte:
		unit = "KB"
		byt = byt / byteKilobyte
	case byt >= byteByte:
		unit = "B"
	default:
		return "See Other..."
	}

	value := strconv.Itoa(byt)

	return value + unit
}
