package templates

import (
	"html/template"
	"io"
	"unicode/utf8"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"

	"hsproom/config"
	"hsproom/plugins"
)

type Template struct {
	Layout   string
	Template string
}

type DefaultMember struct {
	Title string
	User  int
}

func init() {
}
func Del() {
}

func (this *Template) Render(w io.Writer, member interface{}) error {

	return template.Must(template.New(this.Layout).Funcs(map[string]interface{}{
		"linkCSS":    linkCSS,
		"embedImage": embedImage,
		"linkJS":     linkJS,
		"plugin":     plugin,
		"markdown":   markdown,
		"subString":  subString,
	}).ParseFiles(config.LayoutsPath+this.Layout, config.TemplatesPath+this.Template)).Execute(w, member)
}

func linkCSS(cssFile string) template.HTML {
	return template.HTML("<link rel='stylesheet' href='/" + config.CssPath + cssFile + "' type='text/css' />")
}
func embedImage(imgFile string, alt string) template.HTML {
	return template.HTML("<img alt='" + alt + "' src='/" + config.ImgPath + imgFile + "' />")
}
func linkJS(jsFile string) template.HTML {
	return template.HTML("<script type='text/javascript' src='/" + config.JsPath + jsFile + "' ></script>")
}
func plugin(name string) template.HTML {
	return plugins.Plugins[name]()
}
func markdown(markdown string) template.HTML {
	return template.HTML(bluemonday.UGCPolicy().SanitizeBytes(blackfriday.MarkdownCommon([]byte(markdown))))
}
func subString(source string, from int, number int) string {

	var (
		count      int
		total_size int
	)

	if from < 0 {
		from = 0
	}

	if number <= 0 {
		return ""
	}

	for count < from {
		if total_size >= len(source) {
			return ""
		}
		_, size := utf8.DecodeRuneInString(source)
		total_size += size
		count++
	}

	source = source[total_size:]
	count = 0
	total_size = 0

	for total_size < len(source) && count < number {
		_, size := utf8.DecodeRuneInString(source)
		total_size += size
		count++
	}

	return source[:total_size]
}
