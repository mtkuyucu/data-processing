package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type Data struct {
	Id       int       `json:"id"`
	Units    int       `json:"units"`
	Date     time.Time `json:"date"`
	Location string    `json:"location"`
}

func generateRandomData() {
	rand.Seed(time.Now().UnixNano())

	data := Data{
		Id:       rand.Intn(10000),
		Units:    rand.Intn(10000),
		Date:     time.Now().AddDate(0, 0, rand.Intn(365)),
		Location: []string{"PARIS", "LONDON", "NEW YORK", "TOKYO", "SYDNEY"}[rand.Intn(5)],
	}

	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	resp, err := http.Post("http://localhost:8002/streaming-data", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}

	defer resp.Body.Close()

	fmt.Println("Data sent successfully. Status:", resp.Status)
}

func main() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			generateRandomData()
		}
	}
}
