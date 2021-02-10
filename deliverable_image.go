package main

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"strings"
)

type DeliverableImage struct {
	source image.Image
	format string
}

func NewDeliverableImageFromUrl(url string) *DeliverableImage {
	response, httpError := http.Get(url)
	if httpError != nil {
		return nil
	}
	defer response.Body.Close()

	contentType := response.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return nil
	}

	remoteImage, _, decodingError := image.Decode(response.Body)
	if decodingError != nil {
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

	switch format {
	case "jpeg", "jpg":
		_ = jpeg.Encode(buffer, image.source, nil)
	case "png":
		_ = png.Encode(buffer, image.source)
	case "gif":
		_ = gif.Encode(buffer, image.source, nil)
	default:
		buffer = image.encode(image.format)
	}

	return buffer
}
