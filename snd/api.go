package snd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	BASE_URL = "http://api.soundcloud.com"
)

type Client struct {
	Id string
}

// {
// 	"id": 2511,
// 	"kind": "user",
// 	"permalink": "moullinex",
// 	"username": "Moullinex",
// 	"uri": "http://api.soundcloud.com/users/2511",
// 	"permalink_url": "http://soundcloud.com/moullinex",
// 	"avatar_url": "http://i1.sndcdn.com/avatars-000026871864-c10oaq-large.jpg?2aaad5e"
// }
type User struct {
	Id       int64
	Username string
}

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
	// could be "track", ...
	Kind         string
	Id           int64
	UserId       int64
	Title        string
	Description  string
	Uri          string
	Duration     int64
	Commentable  bool
	Streamable   bool
	Downloadable bool
	// could be "all"
	EmbeddableBy   string
	TagList        string
	Permalink      string
	PermalinkUrl   string
	StreamUrl      string
	AttachmentsUri string
	// could be "finished", ...
	State string
	// could be "public"
	Sharing             string
	OriginalContentSize int64
	// could be "mp3"
	OriginalFormat string
	CreatedAt      time.Time

	// statistics
	PlaybackCount    int64
	DownloadCount    int64
	FavoritingsCount int64
	CommentCount     int64

	User User
}

func (c *Client) Resolve(url string) (string, error) {
	uri := fmt.Sprintf("/resolve.json?url=%s&client_id=%s", url, c.Id)

	req, err := http.NewRequest("GET", BASE_URL+uri, nil)
	if err != nil {
		return "", err
	}

	// we use the default roundtripper instead of http.Get() because
	// http.Get() will follow redirects automatically while we want to just
	// fetch the "Location" header of the redirect.
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		switch resp.StatusCode {
		case http.StatusNotFound:
			// if the url was not a valid soundcloud url, we'll get a 404,
			// no need to read the body, just return.
			return "", fmt.Errorf("resource not found")
		default:
			// If the request was not 302 or 404, that's a bit unexpected.
			// not sure what could be the issue now, let's print it out to
			// we learn more.
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return "", fmt.Errorf("couldn't read body: %v", err)
			}
			return "", fmt.Errorf("did not receive a 302 response from /resolve.json, status: %s, body: %s",
				resp.Status, string(body))
		}
	}

	return resp.Header.Get("Location"), nil
}

func (c *Client) GetSound(url string) (*Sound, error) {
	return nil, nil
}
