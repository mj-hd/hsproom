package plugins

var Plugins = map[string]func([]interface{}) interface{}{
	"loadProgram":      loadProgram,
	"byteFormat":       byteFormat,
	"replaceSourceTag": replaceSourceTag,
}

func Init() {

}
func Del() {

}
