package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"github.com/disintegration/imaging"
)

func loadImage(url string) image.Image {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer response.Body.Close()

	remoteImage, _, err := image.Decode(response.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return remoteImage
}

func saveImage(srcImage image.Image) {
	f, err := os.Create("img.jpg")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	jpeg.Encode(f, srcImage, nil)
}

func index(writer http.ResponseWriter, response *http.Request)  {
	url := "https://upload.wikimedia.org/wikipedia/commons/0/0a/Triglav.jpg"
	srcImage := loadImage(url)

	transformedImage := imaging.Resize(srcImage, 128, 128, imaging.Lanczos)

	saveImage(transformedImage)

	fmt.Fprint(writer, "Homepage endpoint hit 1")
}

func handleRequests() {
	http.HandleFunc("/", index)
	log.Fatal(http.ListenAndServe(":80", nil))
}

func main() {
	handleRequests()
}