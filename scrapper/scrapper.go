// scrapper.go: collects all images not already collected from owlrepo's
// screenshot storage and collects them locally (TODO: s3)
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Struct representing relevant data from search_item_listing.json
// see: https://storage.googleapis.com/owlrepo/v1/queries/search_item_listing.json
type SearchItemIndexResponse struct {
	TaskId string `json:"task_id"`
}

func main() {
	// Retrieve the search_item_listing information
	searchItemIndexUrl := "https://storage.googleapis.com/owlrepo/v1/queries/search_item_listing.json"
	searchIndexResults := make([]SearchItemIndexResponse, 0)
	err := getJsonAndDecode(searchItemIndexUrl, &searchIndexResults)
	panicf(err, "Could not retrieve and decode search_item_listing.json")
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
