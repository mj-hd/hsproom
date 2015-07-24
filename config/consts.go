package config

import "time"

const (
	TemplatesPath string = "templates/"
	LayoutsPath   string = "templates/layouts/"
	StaticPath    string = "static/"
	CssPath       string = "static/css/"
	ImgPath       string = "static/img/"
	JsPath        string = "static/js/"

	StartaxLimitSize     int = 1024 * 1024 * 1  // 1MB
	AttachmentsLimitSize int = 1024 * 1024 * 50 // 50MB
	ThumbnailLimitSize   int = 1024 * 1024 * 3  // 3MB
)

// DIRTY:
func JST() *time.Location {
	return time.FixedZone("Asia/Tokyo", 9*60*60)
}
