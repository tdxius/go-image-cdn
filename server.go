package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
)

func homePage(writer http.ResponseWriter, response *http.Request)  {
	url := "https://upload.wikimedia.org/wikipedia/commons/0/0a/Triglav.jpg"
	imageResponse, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer response.Body.Close()

	downloadedImage, _, err := image.Decode(imageResponse.Body)

	f, err := os.Create("img.jpg")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	jpeg.Encode(f, downloadedImage, nil)


	fmt.Fprint(writer, "Homepage endpoint hit 1")
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":80", nil))
}

func main() {
	handleRequests()
}