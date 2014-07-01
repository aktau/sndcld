package main

import (
	"fmt"
	"os"

	"github.com/aktau/sndcld/snd"
)

func main() {
	fmt.Println("Welcome to sndcld. There is no sndcld though, only Zuul.")

	client := snd.Client{Id: "22e566527758690e6feb2b5cb300cc43"}
	url, err := client.Resolve("lol")
	fmt.Println("got:", url, err)
	soundUrl := "https://soundcloud.com/moullinex/discobelle-mix-041-moullinex"
	url, err = client.Resolve(soundUrl)
	if err != nil {
		fmt.Printf("couldn't resolve sound %s: %v\n", soundUrl, err)
		os.Exit(1)
	}
	fmt.Println("got:", url, err)

	sound, err := client.GetSound(url)
	if err != nil {
		fmt.Printf("couldn't fetch sound %s (resolved: %s) metadata: %v\n",
			soundUrl, url, err)
		os.Exit(1)
	}
	fmt.Printf("got sound %s -> %s, %+v\n", soundUrl, url, sound)

	client.DownloadSound(sound)
}
