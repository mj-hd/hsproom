package plugins

var Plugins = map[string]func([]interface{}) interface{}{
	"loadProgram": loadProgram,
	"byteFormat":  byteFormat,
}

func init() {

}
func Del() {

}
