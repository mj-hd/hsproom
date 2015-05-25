package plugins

var Plugins = map[string]func() interface{} {
	"loadProgramHSP3Dish": loadProgramHSP3Dish,
	"loadProgramHGIMG4":   loadProgramHGIMG4,
}

func init() {

}
func Del() {

}
