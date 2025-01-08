// scrapper.go: collects all images not already collected from owlrepo's
// screenshot storage and collects them locally (TODO: s3)
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"scrapper/lvldb"
	"sync"
)

// Struct representing relevant data from search_item_listing.json
// see: https://storage.googleapis.com/owlrepo/v1/queries/search_item_listing.json
type SearchItemIndexResponse struct {
	TaskId string `json:"task_id"`
}

// Struct representing relevant data from slim.json
// see: https://storage.googleapis.com/owlrepo/v1/uploads/97760776-0dfb-4f53-a110-7d1c40e35de0/slim.json
type TaskIdReponse struct {
	TaskId  string
	Url     string
	Payload []struct {
		Screenshot struct {
			Timestamp string `json:"timestamp"`
			FileName  string `json:"name"`
		} `json:"screenshot"`
	} `json:"payload"`
}

func Scrape() {
	db := lvldb.NewLvlDB()
	defer db.Conn.Close()
	err := os.Mkdir("images", os.ModePerm)
	if err != nil && !errors.Is(err, os.ErrExist) {
		fmt.Println(err)
		return
	}
	// Retrieve the search_item_listing information
	searchItemIndexUrl := "https://storage.googleapis.com/owlrepo/v1/queries/search_item_listing.json"
	searchIndexResults := make([]SearchItemIndexResponse, 0)
	err = getJsonAndDecode(searchItemIndexUrl, &searchIndexResults)
	panicf(err, "Could not retrieve and decode search_item_listing.json")

	// For each TaskId retrieve each screenshot's url and attached timestamp
	var wg sync.WaitGroup
	sem := make(chan int, 20)
	searchIndexResultsLen := len(searchIndexResults)
	taskIdErrors := make(chan error, searchIndexResultsLen)
	taskIdResponses := make(chan TaskIdReponse, searchIndexResultsLen)

	for i, sir := range searchIndexResults {
		fmt.Printf("\r processing index search results: [%d / %d]             ", i, searchIndexResultsLen)
		wg.Add(1)
		sem <- 1
		go func(taskId string) {
			defer wg.Done()
			defer func() { <-sem }()

			taskIdUrl := "https://storage.googleapis.com/owlrepo/v1/uploads/" + taskId + "/slim.json"

			taskAlreadyProcessed, err := db.Get(taskIdUrl)
			if err != nil {
				taskIdErrors <- err
				return
			}
			if taskAlreadyProcessed != "" {
				return
			}
			tir := TaskIdReponse{}
			err = getJsonAndDecode(taskIdUrl, &tir)
			if err != nil {
				taskIdErrors <- err
				return
			}
			tir.TaskId = taskId
			tir.Url = taskIdUrl

			taskIdResponses <- tir
		}(sir.TaskId)
	}

	wg.Wait()
	close(taskIdErrors)
	close(taskIdResponses)
	taskIdResponsesLen := len(taskIdResponses)

	if len(taskIdErrors) > (taskIdResponsesLen / 10) {
		err := <-taskIdErrors
		panicf(err, "fetching taskIdResponses had an error rate of > 10%")
	}

	i := 0
	for tir := range taskIdResponses {
		i += 1
		fmt.Printf("\r processing taskIdResponses: [%d / %d]             ", i, taskIdResponsesLen)
		err := db.Put(tir.Url, "processed")
		if err != nil {
			fmt.Println(err)
		}
		for _, payload := range tir.Payload {
			imageUrl := "https://storage.googleapis.com/owlrepo/v1/uploads/" + tir.TaskId + "/raw/" + payload.Screenshot.FileName
			imageKey := tir.TaskId + "~" + payload.Screenshot.Timestamp + "~" + payload.Screenshot.FileName
			wg.Add(1)
			sem <- 1
			go handleImageDownload(db, &wg, sem, imageUrl, imageKey, payload.Screenshot.Timestamp)
		}
	}

	wg.Wait()
}

// TODO: Update this to s3 when done testing.
// Takes a given owl screenshot and downloads it
func handleImageDownload(db lvldb.LvlDB, wg *sync.WaitGroup, sem <-chan int, imageUrl string, imageKey string, timestamp string) {
	defer wg.Done()
	defer func() { <-sem }()
	fileExists, err := db.Get(imageKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	if fileExists != "" {
		return
	}
	resp, err := http.Get(imageUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	file, err := os.Create("./images/" + imageKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = db.Put(imageKey, timestamp)
	if err != nil {
		fmt.Println(err)
	}
}

// Fetches expected json data from url and attempts to decode into target struct
func getJsonAndDecode(url string, target any) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}

// wrapper for critical panic errors
func panicf(err error, message string) {
	if err != nil {
		panic(fmt.Sprintf("%s\n%e", message, err))
	}
}

func main() {
	Scrape()
}
