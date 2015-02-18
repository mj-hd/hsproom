package twitter

type SearchMetadata struct {
	CompletedIn   float32 `json:"completed_in"`
	MaxId         int64   `json:"max_id"`
	MaxIdString   string  `json:"max_id_str"`
	Query         string  `json:"query"`
	RefreshUrl    string  `json:"refresh_url"`
	Count         int     `json:"count"`
	SinceId       int64   `json:"since_id"`
	SinceIdString string  `json:"since_id_str"`
}

type SearchResponse struct {
	Statuses []Tweet        `json:"statuses"`
	Metadata SearchMetadata `json:"search_metadata"`
}

type Tweet struct {
	Contributors         []Contributor     `json:"contributors"` // Not yet generally available to all, so hard to test
	Coordinates          *Coordinates      `json:"coordinates"`
	CreatedAt            string            `json:"created_at"`
	Entities             Entities          `json:"entities"`
	FavoriteCount        int               `json:"favorite_count"`
	Favorited            bool              `json:"favorited"`
	FilterLevel          string            `json:"filter_level"`
	Id                   int64             `json:"id"`
	IdStr                string            `json:"id_str"`
	InReplyToScreenName  string            `json:"in_reply_to_screen_name"`
	InReplyToStatusID    int64             `json:"in_reply_to_status_id"`
	InReplyToStatusIdStr string            `json:"in_reply_to_status_id_str"`
	InReplyToUserID      int64             `json:"in_reply_to_user_id"`
	InReplyToUserIdStr   string            `json:"in_reply_to_user_id_str"`
	Lang                 string            `json:"lang"`
	Place                Place             `json:"place"`
	PossiblySensitive    bool              `json:"possibly_sensitive"`
	RetweetCount         int               `json:"retweet_count"`
	Retweeted            bool              `json:"retweeted"`
	RetweetedStatus      *Tweet            `json:"retweeted_status"`
	Source               string            `json:"source"`
	Scopes               map[string]string `json:"scopes"`
	Text                 string            `json:"text"`
	Truncated            bool              `json:"truncated"`
	User                 User              `json:"user"`
	WithheldCopyright    bool              `json:"withheld_copyright"`
	WithheldInCountries  []string          `json:"withheld_in_countries"`
	WithheldScope        string            `json:"withheld_scope"`

	//Geo is deprecated
	//Geo                  interface{} `json:"geo"`
}
type Contributor struct {
	Id         int64  `json:"id"`
	IdStr      string `json:"id_str"`
	ScreenName string `json:"screen_name"`
}

type Coordinates struct {
	Coordinates [2]float64 `json:"coordinates"` // Coordinate always has to have exactly 2 values
	Type        string     `json:"type"`
}
type User struct {
	ContributorsEnabled            bool     `json:"contributors_enabled"`
	CreatedAt                      string   `json:"created_at"`
	DefaultProfile                 bool     `json:"default_profile"`
	DefaultProfileImage            bool     `json:"default_profile_image"`
	Description                    string   `json:"description"`
	Entities                       Entities `json:"entities"`
	FavouritesCount                int      `json:"favourites_count"`
	FollowRequestSent              bool     `json:"follow_request_sent"`
	FollowersCount                 int      `json:"followers_count"`
	Following                      bool     `json:"following"`
	FriendsCount                   int      `json:"friends_count"`
	GeoEnabled                     bool     `json:"geo_enabled"`
	Id                             int64    `json:"id"`
	IdStr                          string   `json:"id_str"`
	IsTranslator                   bool     `json:"is_translator"`
	Lang                           string   `json:"lang"` // BCP-47 code of user defined language
	ListedCount                    int64    `json:"listed_count"`
	Location                       string   `json:"location"` // User defined location
	Name                           string   `json:"name"`
	Notifications                  bool     `json:"notifications"`
	ProfileBackgroundColor         string   `json:"profile_background_color"`
	ProfileBackgroundImageURL      string   `json:"profile_background_image_url"`
	ProfileBackgroundImageUrlHttps string   `json:"profile_background_image_url_https"`
	ProfileBackgroundTile          bool     `json:"profile_background_tile"`
	ProfileBannerURL               string   `json:"profile_banner_url"`
	ProfileImageURL                string   `json:"profile_image_url"`
	ProfileImageUrlHttps           string   `json:"profile_image_url_https"`
	ProfileLinkColor               string   `json:"profile_link_color"`
	ProfileSidebarBorderColor      string   `json:"profile_sidebar_border_color"`
	ProfileSidebarFillColor        string   `json:"profile_sidebar_fill_color"`
	ProfileTextColor               string   `json:"profile_text_color"`
	ProfileUseBackgroundImage      bool     `json:"profile_use_background_image"`
	Protected                      bool     `json:"protected"`
	ScreenName                     string   `json:"screen_name"`
	ShowAllInlineMedia             bool     `json:"show_all_inline_media"`
	Status                         *Tweet   `json:"status"` // Only included if the user is a friend
	StatusesCount                  int64    `json:"statuses_count"`
	TimeZone                       string   `json:"time_zone"`
	URL                            string   `json:"url"` // From UTC in seconds
	UtcOffset                      int      `json:"utc_offset"`
	Verified                       bool     `json:"verified"`
	WithheldInCountries            string   `json:"withheld_in_countries"`
	WithheldScope                  string   `json:"withheld_scope"`
}
type UrlEntity struct {
	Urls []struct {
		Indices      []int  `json:"indices"`
		Url          string `json:"url"`
		Display_url  string `json:"display_url"`
		Expanded_url string `json:"expanded_url"`
	}
}
type Entities struct {
	Hashtags []struct {
		Indices []int  `json:"indices"`
		Text    string `json:"text"`
	} `json:"hashtags"`
	Urls []struct {
		Indices      []int  `json:"indices"`
		Url          string `json:"url"`
		Display_url  string `json:"display_url"`
		Expanded_url string `json:"expanded_url"`
	} `json:"urls"`
	Url           UrlEntity `json:"url"`
	User_mentions []struct {
		Name        string `json:"name"`
		Indices     []int  `json:"indices"`
		Screen_name string `json:"screen_name"`
		Id          int64  `json:"id"`
		Id_str      string `json:"id_str"`
	} `json:"user_mentions"`
	Media []struct {
		Id              int64      `json:"id"`
		Id_str          string     `json:"id_str"`
		Media_url       string     `json:"media_url"`
		Media_url_https string     `json:"media_url_https"`
		Url             string     `json:"url"`
		Display_url     string     `json:"display_url"`
		Expanded_url    string     `json:"expanded_url"`
		Sizes           MediaSizes `json:"sizes"`
		Type            string     `json:"type"`
		Indices         []int      `json:"indices"`
	}
}
type MediaSizes struct {
	Medium MediaSize `json:"medium"`
	Thumb  MediaSize `json:"thumb"`
	Small  MediaSize `json:"small"`
	Large  MediaSize `json:"large"`
}
type MediaSize struct {
	W      int    `json:"w"`
	H      int    `json:"h"`
	Resize string `json:"resize"`
}
type Place struct {
	Attributes  map[string]string `json:"attributes"`
	BoundingBox struct {
		Coordinates [][][]float64 `json:"coordinates"`
		Type        string        `json:"type"`
	} `json:"bounding_box"`
	ContainedWithin []struct {
		Attributes  map[string]string `json:"attributes"`
		BoundingBox struct {
			Coordinates [][][]float64 `json:"coordinates"`
			Type        string        `json:"type"`
		} `json:"bounding_box"`
		Country     string `json:"country"`
		CountryCode string `json:"country_code"`
		FullName    string `json:"full_name"`
		ID          string `json:"id"`
		Name        string `json:"name"`
		PlaceType   string `json:"place_type"`
		URL         string `json:"url"`
	} `json:"contained_within"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	FullName    string `json:"full_name"`
	Geometry    struct {
		Coordinates [][][]float64 `json:"coordinates"`
		Type        string        `json:"type"`
	} `json:"geometry"`
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	PlaceType string   `json:"place_type"`
	Polylines []string `json:"polylines"`
	URL       string   `json:"url"`
}
