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

var (
	RuntimeVersions map[string]string = map[string]string{
		"hsp3.5b1mod": "HSP 3.5b1(HSP部屋版)",
		"hsp3.5b2mod": "HSP 3.5b2(HSP部屋版)",
		"hsp3.5b3mod": "HSP 3.5b3(HSP部屋版)",
		"hsp3.5b4mod": "HSP 3.5b4(HSP部屋版)",
		"hsp3.5b5mod": "HSP 3.5b5(HSP部屋版)",
	}
)

func IsValidRuntimeVersion(version string) bool {
	for k, _ := range RuntimeVersions {
		if version == k {
			return true
		}
	}

	return false
}

func GetLabelForVersion(version string) string {
	if !IsValidRuntimeVersion(version) {
		return ""
	}

	return RuntimeVersions[version]
}

// DIRTY:
func JST() *time.Location {
	return time.FixedZone("Asia/Tokyo", 9*60*60)
}
