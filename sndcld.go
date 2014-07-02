package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aktau/sndcld/snd"
)

var (
	help = flag.Bool("help", false, "print help")
)

func init() {
	flag.Parse()
}

func usage() {
	fmt.Println("Usage: sndcld [-v] <url-of-track>")
	os.Exit(0)
}

func main() {
	if *help {
		usage()
	}

	var soundUrl string
	if flag.NArg() < 1 {
		// usage()
		// for debugging, we assign a predefined sound
		soundUrl = "https://soundcloud.com/moullinex/discobelle-mix-041-moullinex"
	} else {
		soundUrl = flag.Arg(0)
	}

	fmt.Println("Welcome to sndcld. There is no sndcld though, only Zuul.")

	client := snd.Client{Id: "22e566527758690e6feb2b5cb300cc43"}

	url, err := client.Resolve(soundUrl)
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

	fname, err := client.DownloadSound(sound)
	if err != nil {
		fmt.Println("couldn't download sound:", err)
		os.Exit(1)
	}

	if err := sound.CompleteTags(fname); err != nil {
		fmt.Println("couldn't tag file:", err)
		os.Exit(1)
	}
}
