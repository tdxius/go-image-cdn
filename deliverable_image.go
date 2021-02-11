package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/patrickmn/go-cache"
	"image"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

type DeliverableImage struct {
	source image.Image
	format string
}

var cacher = cache.New(time.Hour, 2*time.Hour)

func cacheOrFetchResponse(url string) *http.Response {
	cachedResponse, found := cacher.Get(url)
	if found {
		reader := bufio.NewReader(bytes.NewReader(cachedResponse.([]byte)))
		response, _ := http.ReadResponse(reader, nil)
		return response
	}

	response, httpError := http.Get(url)
	if httpError != nil {
		return nil
	}
	defer response.Body.Close()

	responseBody, _ := httputil.DumpResponse(response, true)
	cacher.Set(url, responseBody, time.Hour)

	return response
}

func NewDeliverableImageFromUrl(url string) *DeliverableImage {
	response := cacheOrFetchResponse(url)
	contentType := response.Header.Get("Content-Type")

	if !strings.HasPrefix(contentType, "image/") {
		return nil
	}

	remoteImage, _, decodingError := image.Decode(response.Body)
	if decodingError != nil {
		fmt.Println(decodingError)
		return nil
	}

	return &DeliverableImage{
		source: remoteImage,
		format: strings.Split(contentType, "/")[1],
	}
}

func (image DeliverableImage) scale(width int, height int) DeliverableImage {
	if width != 0 && height != 0 {
		image.source = imaging.Fit(image.source, width, height, imaging.Lanczos)
	}

	if width == 0 && height != 0 {
		image.source = imaging.Resize(image.source, 0, height, imaging.Lanczos)
	}

	if height == 0 && width != 0 {
		image.source = imaging.Resize(image.source, width, 0, imaging.Lanczos)
	}

	return image
}

func (image DeliverableImage) encode(format string) *bytes.Buffer {
	buffer := new(bytes.Buffer)

	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		_ = imaging.Encode(buffer, image.source, imaging.JPEG)
	case "png":
		_ = imaging.Encode(buffer, image.source, imaging.PNG)
	case "tiff", "tif":
		_ = imaging.Encode(buffer, image.source, imaging.TIFF)
	case "bmp":
		_ = imaging.Encode(buffer, image.source, imaging.BMP)
	case "gif":
		_ = imaging.Encode(buffer, image.source, imaging.GIF)
	default:
		buffer = image.encode(image.format)
	}

	image.format = strings.ToLower(format)

	return buffer
}
