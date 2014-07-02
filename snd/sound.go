package snd

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mikkyang/id3-go"
)

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

	return artist + " - " + s.NormalizedTitle()
}

var cutExpressions = []*regexp.Regexp{
	regexp.MustCompile(`(?i)[\(\s]*free\s*downloads?[\)\s]*`),
}

// NormalizedTitle - try to clean up the title as returned by the API
func (s *Sound) NormalizedTitle() string {
	title := s.Title

	// sometimes the given title includes the artist's name in the
	// beginning, usually this is wrong (unless it's a self-titled record),
	// so let's avoid that
	if strings.HasPrefix(title, s.User.Username) {
		title = title[len(s.User.Username):]
	}

	title = strings.Replace(title, "/", "-", -1)

	// strip some special characters
	title = stripRunes(title, "*")

	// strip useless annotations in the title (such as "free download")
	for _, re := range cutExpressions {
		title = re.ReplaceAllLiteralString(title, "")
	}

	// now there's possibly some trailing and leading symbols we don't want
	return strings.Trim(title, " \t-")
}

// CompleteTags - writes the Sound's meta information to file fname as tags
//
// Does not overwrite existing information (assumes a field is already
// properly tagged if it exists).
//
// Will return an error if the format is unrecognized or something goes
// wrong while opening the file.
func (s *Sound) CompleteTags(fname string) error {
	if filepath.Ext(fname) != ".mp3" {
		return errors.New("unsupported file format for tagging")
	}

	file, err := id3.Open(fname)
	if err != nil {
		return err
	}
	defer file.Close()

	if prevArtist := file.Artist(); prevArtist != "" {
		fmt.Println(fname, "already had an artist id3 tag:", prevArtist)
	} else {
		fmt.Println(fname, "setting artist:", s.User.Username)
		file.SetArtist(s.User.Username)
	}

	if prevTitle := file.Title(); prevTitle != "" {
		fmt.Println(fname, "already had a title id3 tag:", prevTitle)
	} else {
		fmt.Println(fname, "setting title:", s.NormalizedTitle())
		file.SetTitle(s.NormalizedTitle())
	}

	return nil
}

func stripRunes(str, chr string) string {
	return strings.Map(func(r rune) rune {
		if strings.IndexRune(chr, r) < 0 {
			return r
		}
		return -1
	}, str)
}
