package main

import (
	"github.com/joho/godotenv"
	"github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/memory"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
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
	startTime := time.Now()

	baseUrl := os.Getenv("SOURCE_URL")
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

	duration := int(time.Now().Sub(startTime).Milliseconds())
	fmt.Println(strconv.Itoa(duration) + "ms")
}

func initDotenv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func createCacheClient() *cache.Client {
	memcached, err := memory.NewAdapter(
		memory.AdapterWithAlgorithm(memory.LRU),
		memory.AdapterWithCapacity(10000000),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cacheClient, err := cache.NewClient(
		cache.ClientWithAdapter(memcached),
		cache.ClientWithTTL(10 * time.Minute),
		cache.ClientWithRefreshKey("opn"),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return cacheClient
}

func main() {
	initDotenv()
	cacheClient := createCacheClient()

	handler := http.HandlerFunc(index)

	http.Handle("/", cacheClient.Middleware(handler))

	fmt.Println("Web server is listening on port 80")
	serverError := http.ListenAndServe(":80", nil)
	if serverError != nil {
		fmt.Println("Failed to start web server on port 80.")
	}
}
