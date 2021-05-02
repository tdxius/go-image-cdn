package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"strconv"
	"strings"
)

func main() {
	initDotenv()

	router := gin.Default()
	router.GET("/*path", index)
	router.Run()
}

func index(context *gin.Context) {
	imageUrl := strings.TrimLeft(context.Request.URL.Path, "/")
	fmt.Println(imageUrl)

	deliverableImage := NewDeliverableImageFromUrl(imageUrl)
	if deliverableImage == nil {
		fmt.Println("No image found at URL: " + imageUrl)
		context.AbortWithStatus(404)
		return
	}

	width := convertStringToInt(context.Query("width"))
	height := convertStringToInt(context.Query("height"))
	format := context.Query("format")
	if format == "" {
		format = deliverableImage.format
	}

	bufferedImage := deliverableImage.scale(width, height).encode(format)

	imageResponse(context, bufferedImage, format)
}

func initDotenv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func convertStringToInt(string string) int {
	value, _ := strconv.ParseInt(string, 10, 16)
	return int(value)
}

func imageResponse(context *gin.Context, buffer *bytes.Buffer, format string) {
	context.Writer.WriteHeader(200)
	context.Writer.Header().Set("Content-Type", "image/" + format)
	context.Writer.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))

	_, writeError := context.Writer.Write(buffer.Bytes())
	if writeError != nil {
		log.Println("unable to write source.")
	}
	context.Writer.WriteHeaderNow()
}
