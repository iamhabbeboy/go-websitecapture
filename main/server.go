package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type thumbnailRequest struct {
	Url string `json:"url"`
}

type screenshotAPIRequest struct {
	Url            string `json:"url"`
	Output         string `json:"output"`
	Width          int    `json:"width"`
	Height         int    `json:"height"`
	ThumbnailWidth int    `json:"thumbnail_width"`
}

type response struct {
	screenshot string `json:"screenshot"`
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	http.HandleFunc("/api/thumbnail", thumbnailHandler)

	fmt.Println("Server listening on port 3000")
	log.Panic(http.ListenAndServe(":3000", nil))
}

func thumbnailHandler(w http.ResponseWriter, r *http.Request) {
	var decoded thumbnailRequest

	err := json.NewDecoder(r.Body).Decode(&decoded)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	apiRequest := screenshotAPIRequest{
		Url:            decoded.Url,
		Output:         "json",
		Width:          1920,
		Height:         1080,
		ThumbnailWidth: 300,
	}

	resultChan := make(chan response, 1)
	go processWebsiteThumbnail(apiRequest, resultChan)
	res := <-resultChan
	close(resultChan)

	if err != nil {
		log.Fatalf("Error occured %s", err)
	}
	fmt.Println(string(res.screenshot))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res.screenshot))

	fmt.Println("Processing job ....🔥")
	fmt.Printf("Got the following url: %s\n", decoded.Url)
}

/*
 * A go routine that fetches data from an endpoint
 */
func processWebsiteThumbnail(apiRequest screenshotAPIRequest, resultChan chan response) {
	apiUrl := os.Getenv("API_URL")
	apiToken := os.Getenv("API_TOKEN")

	resp, err := http.Get(apiUrl + "/screenshot?token=" + apiToken + "&url=" + apiRequest.Url)
	if err != nil {
		log.Fatal("Error occured while getting image..", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	res := response{screenshot: string(body)}
	resultChan <- res
}
