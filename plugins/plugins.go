package plugins

var Plugins = map[string]func() interface{}{
	"LoadProgramScript": loadProgramScript,
}

func init() {

}
func Del() {

}
