package snd

import "strings"

// {
// 	"kind": "track",
// 	"id": 152514285,
// 	"created_at": "2014/06/02 22:04:29 +0000",
// 	"user_id": 2511,
// 	"duration": 4199289,
// 	"commentable": true,
// 	"state": "finished",
// 	"original_content_size": 167959022,
// 	"sharing": "public",
// 	"tag_list": "Gomma Discotexas \"Love Magnetic\"",
// 	"permalink": "discobelle-mix-041-moullinex",
// 	"streamable": true,
// 	"embeddable_by": "all",
// 	"downloadable": false,
// 	"purchase_url": null,
// 	"label_id": null,
// 	"purchase_title": null,
// 	"genre": "Mixtape",
// 	"title": "Discobelle Mix 041: Moullinex",
// 	"description": "A new mixtape, long overdue!\r\n'Love Magnetic' EP out June 13 on Gomma Records",
// 	"label_name": "",
// 	"release": "",
// 	"track_type": "podcast",
// 	"key_signature": "",
// 	"isrc": "",
// 	"video_url": null,
// 	"bpm": null,
// 	"release_year": null,
// 	"release_month": null,
// 	"release_day": null,
// 	"original_format": "mp3",
// 	"license": "all-rights-reserved",
// 	"uri": "http://api.soundcloud.com/tracks/152514285",
// 	"user": {},
// 	"permalink_url": "http://soundcloud.com/moullinex/discobelle-mix-041-moullinex",
// 	"artwork_url": "http://i1.sndcdn.com/artworks-000081264242-cftxc0-large.jpg?2aaad5e",
// 	"waveform_url": "http://w1.sndcdn.com/WYAiN8pZW7Bv_m.png",
// 	"stream_url": "http://api.soundcloud.com/tracks/152514285/stream",
// 	"playback_count": 24556,
// 	"download_count": 0,
// 	"favoritings_count": 1387,
// 	"comment_count": 128,
// 	"attachments_uri": "http://api.soundcloud.com/tracks/152514285/attachments"
// }
type Sound struct {
	Id                  int64  `json:"id"`
	Kind                string `json:"kind"` // could be "track", ...
	UserId              int64  `json:"user_id"`
	Title               string `json:"title"`
	Description         string `json:"description"`
	Uri                 string `json:"uri"`
	Duration            int64  `json:"duration"`
	Commentable         bool   `json:"commentable"`
	Streamable          bool   `json:"streamable"`
	Downloadable        bool   `json:"downloadable"`
	EmbeddableBy        string `json:"embeddable_by"` // could be "all"
	TagList             string `json:"tag_list"`
	Permalink           string `json:"permalink"`
	PermalinkUrl        string `json:"permalink_url"`
	StreamUrl           string `json:"stream_url"`
	DownloadUrl         string `json:"download_url"`
	AttachmentsUri      string `json:"attachments_uri"`
	State               string `json:"state"`   // could be "finished", ...
	Sharing             string `json:"sharing"` // could be "public"
	OriginalContentSize int64  `json:"original_content_size"`
	OriginalFormat      string `json:"original_format"` // could be "mp3"
	CreatedAt           Time   `json:"created_at"`

	// statistics
	PlaybackCount    int64 `json:"playback_count"`
	DownloadCount    int64 `json:"download_count"`
	FavoritingsCount int64 `json:"favoritings_count"`
	CommentCount     int64 `json:"comment_count"`

	// user sub-object
	User User `json:"user"`
}

func (s *Sound) Filename() string {
	// TODO: the user is not always the artist, find a better heuristic...
	artist := strings.Replace(s.User.Username, "/", "-", -1)
	title := strings.Replace(s.Title, "/", "-", -1)
	return artist + " - " + title + ".mp3"
}
