package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"github.com/disintegration/imaging"
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

func encodeImage(srcImage image.Image, format string) *bytes.Buffer {
	buffer := new(bytes.Buffer)
	var encodingError error

	switch format {
	case "jpeg":
		encodingError = jpeg.Encode(buffer, srcImage, nil)
	case "png":
		encodingError = png.Encode(buffer, srcImage)
	default:
		encodingError = jpeg.Encode(buffer, srcImage, nil)
	}

	if encodingError != nil {
		fmt.Println("Unable to encode image: " + encodingError.Error())
	}

	return buffer
}

func scaleImage(srcImage image.Image, width int, height int) image.Image {
	if width != 0 && height != 0 {
		return imaging.Fit(srcImage, width, height, imaging.Lanczos)
	}

	if width == 0 && height != 0 {
		return imaging.Resize(srcImage, 0, height, imaging.Lanczos)
	}

	if height == 0 && width != 0 {
		return imaging.Resize(srcImage, width, 0, imaging.Lanczos)
	}

	return srcImage
}

func writeImage(writer http.ResponseWriter, buffer *bytes.Buffer) {
	writer.Header().Set("Content-Type", "image/jpeg")
	writer.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))

	_, writeError := writer.Write(buffer.Bytes())
	if writeError != nil {
		log.Println("unable to write image.")
	}
}

func index(writer http.ResponseWriter, response *http.Request) {
	baseUrl := "https://enavtika.si"
	imageUrl := baseUrl + response.URL.Path
	fmt.Println(imageUrl)

	srcImage := loadImage(imageUrl)
	if srcImage == nil {
		fmt.Fprint(writer, "No image found at URL: "+imageUrl)
	}

	width, _ := strconv.ParseInt(response.URL.Query().Get("width"), 10, 16)
	height, _ := strconv.ParseInt(response.URL.Query().Get("height"), 10, 16)
	scaledImage := scaleImage(srcImage, int(width), int(height))

	format := response.URL.Query().Get("format")
	encodedImage := encodeImage(scaledImage, format)

	writeImage(writer, encodedImage)
}

func main() {
	http.HandleFunc("/", index)
	fmt.Println("Web server is listening on port 80")
	serverError := http.ListenAndServe(":80", nil)
	if serverError != nil {
		fmt.Println("Failed to start web server on port 80.")
	}
}
