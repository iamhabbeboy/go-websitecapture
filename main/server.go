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
	go testing(apiRequest, resultChan)
	res := <-resultChan
	fmt.Printf(res.screenshot)
	close(resultChan)
	// go func (apiRequest, w, screenshotAPIResponse) {
	// go processWebsiteThumbnail(apiRequest, w, screenshotAPIResponse)
	// 	if err != nil {
	// 		log.Fatal("Error occured while getting data")
	// 	}
	// 	fmt.Printf(string(screenshot))
	// }

	// jsonResp, err := json.Marshal(response)

	// fmt.Printf(string(jsonResp))

	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// w.Write(jsonResp)

	// if err != nil {
	// 	log.Fatalf("Error occured %s", err)
	// }
	fmt.Printf("Got the following url: %s\n", decoded.Url)
}

func testing(apiRequest screenshotAPIRequest, resultChan chan response) {
	// name := apiRequest.Url
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
	res := json.Unmarshal(body, &response{})
	fmt.Println(string(body))
	// res := response{screenshot: name}
	resultChan <- string(res)
}

func checkError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// func processWebsiteThumbnail(apiRequest screenshotAPIRequest, w http.ResponseWriter, screenshot screenshotAPIResponse) {

// 	apiUrl := os.Getenv("API_URL")
// 	apiToken := os.Getenv("API_TOKEN")

// 	fmt.Printf(apiUrl)

// 	response, err := http.Get(apiUrl + "/screenshot?token=" + apiToken + "&url=" + apiRequest.Url)

// 	if err != nil {
// 		return nil, err
// 	}

// 	fmt.Println("Processing job ....🔥")

// 	err = json.NewDecoder(response.Body).Decode(&screenshot)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return screenshot, nil
// }
