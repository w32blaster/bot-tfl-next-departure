package commands

// this file contains functions working with TFL

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"time"

	"net/http"
	"net/url"

	"gopkg.in/telegram-bot-api.v4"
)

type TFLInlineQueryResult struct {
	Total   int       `json:"total"`
	Query   string    `json:"query"`
	Matches []Station `json:"matches"`
}

type Station struct {
	IcsID string   `json:"icsId"`
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Modes []string `json:"modes"`
}

var httpClient = http.Client{Timeout: time.Duration(3 * time.Second)}

// GetStationListByPattern make a HTTP request to TFL and build up the response
func GetStationListByPattern(searchingPattern string) []interface{} {
	var answers []interface{}

	// make HTTP request to TFL
	results, error := performHTTPRequestToTFLForStations(searchingPattern)
	if error != nil {
		log.Println(error)
		return make([]interface{}, 0)
	}

	// build array of Inline Query results (so called "articles")
	for _, station := range *results {

		// Build one line for inline answer (one result)
		strStationID := fmt.Sprint(station.IcsID)
		answer := tgbotapi.NewInlineQueryResultArticleMarkdown(strStationID, station.Name, "/station-"+strStationID)
		answer.Description = html.EscapeString(printModesInMarkdown(station.Modes))

		answers = append(answers, answer)
	}

	return answers
}

// Call TFL API for a stations
//
// Here we are calling so called "Stop Points", please refer to official documentation:
//      https://api.tfl.gov.uk/swagger/ui/index.html#!/StopPoint/StopPoint_Search
//
// or, here is the example request:
//      https://api.tfl.gov.uk/StopPoint/Search?query=Camden%20Town
//
func performHTTPRequestToTFLForStations(searchingPattern string) (*[]Station, error) {

	// firstly, prepare request URL
	apiURL := "https://api.tfl.gov.uk/StopPoint/Search?query=" + url.PathEscape(searchingPattern) + "&maxResults=10&app_id=4c754c2a&app_key=9eec9fd4bb56bf3732b2627b391d05b9"

	// call API
	resp, err := httpClient.Get(apiURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// parse JSON
	var result TFLInlineQueryResult
	json.NewDecoder(resp.Body).Decode(&result)

	return &result.Matches, nil
}

// simply prints array of strings in markdown syntax
func printModesInMarkdown(arr []string) string {
	if arr == nil {
		return ""
	}

	var buffer bytes.Buffer
	for _, e := range arr {
		buffer.WriteString("● ")
		buffer.WriteString(e)
		buffer.WriteString("\n")
	}

	return buffer.String()
}
