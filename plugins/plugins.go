package plugins

var Plugins = map[string]func() interface{} {
	"loadProgram": loadProgram,
}

func init() {

}
func Del() {

}
