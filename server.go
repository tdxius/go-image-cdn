package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func imageResponse(writer http.ResponseWriter, buffer *bytes.Buffer) {
	writer.Header().Set("Content-Type", "image/jpeg")
	writer.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))

	_, writeError := writer.Write(buffer.Bytes())
	if writeError != nil {
		log.Println("unable to write source.")
	}
}

func getQueryParamAsInt(request *http.Request, key string) int {
	value, _ := strconv.ParseInt(request.URL.Query().Get(key), 10, 16)
	return int(value)
}

func getQueryParams(request *http.Request) (width int, height int, format string) {
	width = getQueryParamAsInt(request, "width")
	height = getQueryParamAsInt(request, "height")
	format = request.URL.Query().Get("format")
	return
}

func index(writer http.ResponseWriter, request *http.Request) {
	baseUrl := "https://enavtika.si"
	imageUrl := baseUrl + request.URL.Path
	fmt.Println(imageUrl)

	deliverableImage := NewDeliverableImageFromUrl(imageUrl)
	if deliverableImage == nil {
		fmt.Println("No image found at URL: " + imageUrl)
		writer.WriteHeader(404)
		return
	}

	width, height, format := getQueryParams(request)
	bufferedImage := deliverableImage.scale(width, height).encode(format)

	imageResponse(writer, bufferedImage)
}

func main() {
	http.HandleFunc("/", index)

	fmt.Println("Web server is listening on port 80")
	serverError := http.ListenAndServe(":80", nil)
	if serverError != nil {
		fmt.Println("Failed to start web server on port 80.")
	}
}
