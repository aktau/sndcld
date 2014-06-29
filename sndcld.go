package main

import (
	"fmt"

	"github.com/aktau/sndcld/snd"
)

func main() {
	fmt.Println("Welcome to sndcld. There is no sndcld though, only Zuul.")

	client := snd.Client{Id: "22e566527758690e6feb2b5cb300cc43"}
	url, err := client.Resolve("lol")
	fmt.Println("got:", url, err)
	url, err = client.Resolve("https://soundcloud.com/moullinex/discobelle-mix-041-moullinex")
	fmt.Println("got:", url, err)
}
