package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	//"github.com/disintegration/imaging"
	"strconv"
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

func writeImage(writer http.ResponseWriter, srcImage image.Image) {
	buffer := new(bytes.Buffer)
	encodeError := jpeg.Encode(buffer, srcImage, nil)
	if encodeError != nil {
		log.Println("unable to encode image.")
	}

	writer.Header().Set("Content-Type", "image/jpeg")
	writer.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))

	_, writeError := writer.Write(buffer.Bytes())
	if writeError != nil {
		log.Println("unable to write image.")
	}
}

func index(writer http.ResponseWriter, response *http.Request)  {
	baseUrl := "https://enavtika.si"
	imageUrl := baseUrl + response.URL.Path
	fmt.Println(imageUrl)

	srcImage := loadImage(imageUrl)
	if srcImage == nil {
		fmt.Fprint(writer, "No image found at URL: " + imageUrl)
	}

	//transformedImage := imaging.Resize(srcImage, 128, 128, imaging.Lanczos)

	writeImage(writer, srcImage)

	fmt.Fprint(writer, response.URL.Path)
}

func handleRequests() {
	http.HandleFunc("/", index)
	log.Fatal(http.ListenAndServe(":80", nil))
}

func main() {
	handleRequests()
}