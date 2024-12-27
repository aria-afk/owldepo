// scrapper.go: collects all images not already collected from owlrepo's
// screenshot storage and collects them locally (TODO: s3)
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	Payload []struct {
		Screenshot struct {
			Timestamp string `json:"timestamp"`
			FileName  string `json:"name"`
		} `json:"screenshot"`
	} `json:"payload"`
}

func Scrape() {
	// Retrieve the search_item_listing information
	searchItemIndexUrl := "https://storage.googleapis.com/owlrepo/v1/queries/search_item_listing.json"
	searchIndexResults := make([]SearchItemIndexResponse, 0)
	err := getJsonAndDecode(searchItemIndexUrl, &searchIndexResults)
	panicf(err, "Could not retrieve and decode search_item_listing.json")

	// For each TaskId retrieve each screenshot's url and attached timestamp
	var wg sync.WaitGroup
	sem := make(chan int, 20)
	taskIdErrors := make(chan error, len(searchIndexResults))
	taskIdResponses := make(chan TaskIdReponse, len(searchIndexResults))

	for i, sir := range searchIndexResults {
		// TESTING: REMOVE ME
		if i > 5 {
			break
		}
		// -----------------
		wg.Add(1)
		sem <- 1
		go func(taskId string) {
			defer wg.Done()
			defer func() { <-sem }()

			taskIdUrl := "https://storage.googleapis.com/owlrepo/v1/uploads/" + taskId + "/slim.json"
			tir := TaskIdReponse{}
			err := getJsonAndDecode(taskIdUrl, &tir)
			if err != nil {
				taskIdErrors <- err
				return
			}

			taskIdResponses <- tir
		}(sir.TaskId)
	}

	wg.Wait()
	close(taskIdErrors)
	close(taskIdResponses)

	if len(taskIdErrors) > (len(taskIdResponses) / 10) {
		err := <-taskIdErrors
		panicf(err, "fetching taskIdResponses had an error rate of > 10%")
	}

	for tir := range taskIdResponses {
		fmt.Println(tir)
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
}
