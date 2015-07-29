package plugins

var Plugins = map[string]func([]interface{}) interface{}{
	"loadProgram":      loadProgram,
	"byteFormat":       byteFormat,
	"replaceSourceTag": replaceSourceTag,
	"showMarkdownHelp": showMarkdownHelp,
	"shorten":          shorten,
}

func Init() {

}
func Del() {

}
