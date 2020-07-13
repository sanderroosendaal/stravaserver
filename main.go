package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Update struct {
	Title string `json:"title"`
}

type Message struct {
	AspectType     string `json:"aspect_type"`
	EventTime      int64  `json:"event_time"`
	ObjectId       int64  `json:"object_id"`
	ObjectType     string `json:"object_type"`
	OwnerId        int64  `json:"owner_id"`
	SubscriptionId int64  `json:"subscription_id"`
	Updates        Update `json:"updates"`
}

// HTTPClient interface
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	// Client to wrap around real client and allow mocks
	Client HTTPClient
)

func init() {
	Client = &http.Client{}
}

func main() {
	update := Update{
		Title: "New Title",
	}

	message := Message{
		AspectType:     "create",
		EventTime:      1516126040,
		ObjectId:       3591534799,
		ObjectType:     "activity",
		OwnerId:        47155909,
		SubscriptionId: 160797,
		Updates:        update,
	}

	requestBody, err := json.Marshal(message)
	fmt.Printf("%s\n", requestBody)

	if err != nil {
		log.Fatalln(err)
		return
	}
	address := "http://localhost:8000/rowers/strava/webhooks/"
	req, err := http.NewRequest("POST", address, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := Client.Do(req)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode != 200 {
		fmt.Printf("Response status code %v\n", resp.StatusCode)
		return
	}

	fmt.Println("Ping OK")
}
