package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
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
	format := response.URL.Query().Get("format")
	fmt.Println(imageUrl + ", Format: " + format)

	srcImage := loadImage(imageUrl)
	if srcImage == nil {
		fmt.Fprint(writer, "No image found at URL: "+imageUrl)
	}

	bufferedImage := encodeImage(srcImage, "png")

	//transformedImage := imaging.Resize(srcImage, 128, 128, imaging.Lanczos)

	writeImage(writer, bufferedImage)

	fmt.Fprint(writer, response.URL.Path)
}

func main() {
	http.HandleFunc("/", index)
	fmt.Println("Web server is listening on port 80")
	serverError := http.ListenAndServe(":80", nil)
	if serverError != nil {
		fmt.Println("Failed to start web server on port 80.")
	}
}
