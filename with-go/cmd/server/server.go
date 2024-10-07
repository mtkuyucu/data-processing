package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Data struct {
	Id       int       `json:"id"`
	Units    int       `json:"units"`
	Date     time.Time `json:"date"`
	Location string    `json:"location"`
}

func handleStreamingData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var data Data
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}
	data = cleanData(data)
	err = saveDataToFile(data)
	if err != nil {
		log.Printf("Error saving data to file: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data received successfully"))
}

func cleanData(data Data) Data {
	data.Location = strings.ToLower(data.Location)
	return data
}

func saveDataToFile(data Data) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %v", err)
	}

	tmpDir := filepath.Join(currentDir, "..", "tmp")

	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return fmt.Errorf("error creating tmp directory: %v", err)
	}

	fileName := filepath.Join(tmpDir, fmt.Sprintf("%s.txt", time.Now().Format("2006-01-02-15-04-05")))

	content := fmt.Sprintf("Location: %s\nUnits: %d\nDate: %s\n\n",
		data.Location,
		data.Units,
		data.Date.Format("2006-01-02"))

	err = ioutil.WriteFile(fileName, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}

func main() {
	http.HandleFunc("/streaming-data", handleStreamingData)

	fmt.Println("Server is running on http://localhost:8002")
	log.Fatal(http.ListenAndServe(":8002", nil))
}
