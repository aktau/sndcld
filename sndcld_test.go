package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/aktau/sndcld/snd"
)

var client = snd.Client{Id: "22e566527758690e6feb2b5cb300cc43"}

func TestRetag(t *testing.T) {
	// download a file (if necessary and tag it multiple times)
	table := []struct {
		inUrl   string
		inFname string
	}{
		{
			inUrl:   "https://soundcloud.com/purple_oslo/sleepers",
			inFname: "PURPURR PURPLE OSLO - Sleepers.mp3",
		},
		{
			inUrl:   "https://soundcloud.com/justinfaustmusic/justin-faust-la-fete-free",
			inFname: "Justin Faust - La Fête.mp3",
		},
	}

	for _, v := range table {
		fmt.Println("TESTING:", v.inFname, v.inUrl)

		// get the sound
		url, err := client.Resolve(v.inUrl)
		if err != nil {
			t.Errorf("couldn't resolve sound %s: %v\n", v.inUrl, err)
		}

		sound, err := client.GetSound(url)
		if err != nil {
			t.Errorf("couldn't fetch sound %s (resolved: %s) metadata: %v\n",
				v.inUrl, url, err)
		}

		// download if fname does not exist
		if _, err := os.Stat(v.inFname); os.IsNotExist(err) {
			v.inFname = download(&client, sound)
			// t.Errorf("no such file or directory: %s", filename)
		}

		for i := 0; i < 5; i++ {
			if err := sound.CompleteTags(v.inFname); err != nil {
				t.Errorf("couldn't tag file on attempt %d: %v", i, err)
			}
		}
	}
}
