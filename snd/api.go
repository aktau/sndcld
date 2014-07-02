package snd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
)

const (
	BASE_URL = "http://api.soundcloud.com"
)

type Client struct {
	Id string
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
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	var sound Sound
	if err := dec.Decode(&sound); err != nil {
		return nil, err
	}

	return &sound, nil
}

func (c *Client) DownloadSound(sound *Sound) (*http.Response, error) {
	fmt.Println("Attempting to download", sound.Title)

	var u string
	if sound.Downloadable {
		fmt.Println("\tdownload version supported, this is probably the highest-quality version", sound.DownloadUrl)
		u = sound.DownloadUrl
	} else {
		fmt.Println("\tno download version, streaming instead:", sound.StreamUrl)
		u = sound.StreamUrl
	}

	v := url.Values{}
	v.Add("client_id", c.Id)
	u = u + "?" + v.Encode()

	fmt.Println("GET -> ", u)
	resp, err := http.Get(u)
	return resp, err
}

// if you have the codec, prefer that to dermine the extension, because it's
// less ambiguous (MIME audio/ogg can be both Vorbis and Opus).
var codecToExt = map[string]string{
	"mp3":    ".mp3",
	"flac":   ".flac",
	"wave":   ".wav",
	"opus":   ".opus",
	"vorbis": ".ogg",
}

// if you don't have a codec, you can use the MIME type to guess an
// extension
var mimeToExt = map[string]string{
	"audio/x-wav": ".wav",
	"audio/mpeg":  ".mp3",
	"audio/mpeg3": ".mp3",
	"audio/mp3":   ".mp3",
	"audio/flac":  ".flac",
	"audio/ogg":   ".ogg",
}

// ParseFiletype - return a MIME type and suggest an extension
//
// Based on the values usually delivered via the Content-Type and
// Content-Disposition HTTP headers. Note: the extension contains the
// starting dot.
func ParseFiletype(contentType, contentDisposition string) (mimeType string, ext string, err error) {
	defer func() {
		// return an error if either the mime type or the extension couldn't
		// be determined after a best effort
		if mimeType == "" || ext == "" {
			err = fmt.Errorf("could not derive MIME type and/or extension from %v and %v",
				contentType, contentDisposition)
		}
	}()

	mimeType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		// couldn't parse, will have to derive from the file extension
		fmt.Println("could not parse Content-Type:", contentType)
	}
	codec := params["codecs"]
	fmt.Printf("Header specified: MIME = '%v', codec = '%v'\n", mimeType, codec)

	// copy the extension of the original file to the new file (later we can
	// do MIME-type based deduction should there be format with multiple
	// extensions and we want to "normalize").
	_, params, err = mime.ParseMediaType(contentDisposition)
	if err != nil {
		fmt.Println("could not parse Content-Disposition", contentDisposition)
	}
	origName := params["filename"]
	fmt.Printf("Header specified: orig. name: '%v'\n", origName)

	ext = determineExtension(origName, codec, mimeType)

	// if we couldn't parse or find the mimetype, decide it based on the
	// file extension (if possible)
	if mimeType == "" && ext != "" {
		guessMime := mime.TypeByExtension(ext)
		mimeType, _, _ = mime.ParseMediaType(guessMime)
	}

	return
}

// determineExtension - try to determine the best extension of the file
//
// The decision tree:
// 1. the extension of the file as suggested by soundcloud (as found in
//    the Content-Disposition header)
// 2. the extension usually associated with the codec (if the codec was
//    found)
// 3. the extension usually associated with the MIME type (of the MIME
//    type was found)
func determineExtension(origName, codec, mimeType string) string {
	if ext := filepath.Ext(origName); ext != "" {
		return ext
	}

	if ext, ok := codecToExt[codec]; ok {
		return ext
	}

	if ext, ok := mimeToExt[mimeType]; ok {
		return ext
	}

	return ""
}
