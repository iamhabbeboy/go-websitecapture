package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type thumbnailRequest struct {
	Url string `json:"url"`
}

type screenshotAPIRequest struct {
	Token          string `json:"token"`
	Url            string `json:"url"`
	Output         string `json:"output"`
	Width          int    `json:"width"`
	Height         int    `json:"height"`
	ThumbnailWidth int    `json:"thumbnail_width"`
}

type screenshotAPIResponse struct {
	screenshot string `json"screenshot"`
}

func main() {
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
		Token:          "6FK6XW3-K4HMDYD-HWNCT73-MV8M0EH",
		Url:            decoded.Url,
		Output:         "json",
		Width:          1920,
		Height:         1080,
		ThumbnailWidth: 300,
	}

	go processWebsiteThumbnail(apiRequest, w)
	fmt.Printf("Got the following url: %s\n", decoded.Url)
}

func checkError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func processWebsiteThumbnail(apiRequest screenshotAPIRequest, w http.ResponseWriter) {
	req, err := http.NewRequest("GET", "https://shot.screenshotapi.net/screenshot?token="+apiRequest.Token+"&url="+apiRequest.Url, nil)
	req.Header.Set("Content-Type", "application/json")

	response, err := (&http.Client{}).Do(req)
	fmt.Println("Processing job ....🔥")
	checkError(err)

	defer response.Body.Close()
	var body []byte
	body, err = ioutil.ReadAll(response.Body)
	checkError(err)

	json.NewEncoder(w).Encode(body)

	_, err = fmt.Fprintf(w, `{ "screenshot": "%s" }`, body)
}
