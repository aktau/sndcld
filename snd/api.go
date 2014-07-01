package snd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
