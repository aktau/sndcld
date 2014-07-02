package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aktau/sndcld/snd"
	"github.com/cheggaaa/pb"
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

	resp, err := client.DownloadSound(sound)
	if err != nil {
		fmt.Println("couldn't download sound:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	fname, err := storeSound(sound, resp)
	if err != nil {
		fmt.Println("couldn't store sound locally:", err)
		os.Exit(1)
	}

	if err := sound.CompleteTags(fname); err != nil {
		fmt.Println("couldn't tag file:", err)
		os.Exit(1)
	}
}

func storeSound(sound *snd.Sound, resp *http.Response) (string, error) {
	h := resp.Header
	_, ext, err := snd.ParseFiletype(h.Get("Content-Type"), h.Get("Content-Disposition"))
	if err != nil {
		// if there's still nothing, assume it's mp3
		ext = ".mp3"
	}

	fname := sound.Filename() + ext

	f, err := os.Create(fname)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// create progressbar
	nbytes, _ := strconv.Atoi(h.Get("Content-Length"))

	bar := pb.New(nbytes).SetUnits(pb.U_BYTES).SetRefreshRate(time.Millisecond * 10)
	bar.ShowSpeed = true
	bar.Start()
	defer bar.Finish()

	// create multi writer, write to the progress bar and the file at the
	// same time.
	writer := io.MultiWriter(f, bar)

	// stream the http response to the file (check if chunked encoding
	// doesn't mess with this)
	_, err = io.Copy(writer, resp.Body)
	return fname, err
}
