package plugins

var Plugins = map[string]func([]interface{}) interface{}{
	"loadProgram": loadProgram,
}

func init() {

}
func Del() {

}
