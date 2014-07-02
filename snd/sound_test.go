package snd

import "testing"

func TestNormalizedTitle(t *testing.T) {
	table := []struct {
		in  Sound
		out string
	}{
		{in: Sound{Title: "My song ( frEe  downLOAD)"}, out: "My song"},
		{in: Sound{Title: "       free download   "}, out: ""},
		{in: Sound{Title: "free downloads for life"}, out: "for life"},
		{in: Sound{Title: "Artist - Song", User: User{Username: "Artist"}}, out: "Song"},
	}

	for _, input := range table {
		if input.in.NormalizedTitle() != input.out {
			t.Error(input.in.NormalizedTitle(), "should've been", input.out)
		}
	}
}
