package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"text/template"
	"time"
)

type StatusData struct {
	Status      Stats  `json:"status"`
	WaterStatus string `json:"-"`
	WindStatus  string `json:"-"`
}

type Stats struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

var (
	data StatusData
)

func main() {
	go startService()
	http.HandleFunc("/", mainHandler)
	http.ListenAndServe(":8080", nil)
}

func startService() {
	ticker := time.NewTicker(4 * time.Second)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				fmt.Println("At ", t)
				doService()
				fmt.Println("=======================================================")
			}
		}

	}()

	time.Sleep(60 * time.Minute)
	ticker.Stop()
	done <- true
	fmt.Println("Stopped service")
}

func doService() {
	file, err := os.ReadFile("status.json")
	if err != nil {
		fmt.Println("Error reading status.json: ", err)
		return
	}

	err = json.Unmarshal(file, &data)
	if err != nil {
		fmt.Println("Error unmarshaling status.json: ", err)
		return
	}
	data.Status.Water = GenerateRandomNumber()
	data.Status.Wind = GenerateRandomNumber()
	data.WaterStatus, data.WindStatus = SetStatus(data.Status)

	newFile, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error marshal data to json: ", err)
		return
	}

	err = os.WriteFile("status.json", newFile, 0644)
	if err != nil {
		fmt.Println("Error writing data to files status.json: ", err)
		return
	}

	fmt.Println("Water :", data.Status.Water, "Status", data.WaterStatus)
	fmt.Println("Wind  :", data.Status.Wind, "Status", data.WindStatus)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	tmplt, err := template.ParseFiles("template.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmplt.Execute(w, data)
}

func SetStatus(s Stats) (string, string) {
	var waterStatus string
	var windStatus string

	//waterStatus
	if s.Water <= 5 {
		waterStatus = "Aman"
	} else if s.Water >= 6 && s.Water <= 8 {
		waterStatus = "Siaga"
	} else {
		waterStatus = "Bahaya"
	}

	//windStatus
	if s.Wind <= 6 {
		windStatus = "Aman"
	} else if s.Wind >= 7 && s.Wind <= 15 {
		windStatus = "Siaga"
	} else {
		windStatus = "Bahaya"
	}

	return waterStatus, windStatus
}

func GenerateRandomNumber() int {
	min := 1
	max := 100

	rand.Seed(time.Now().UnixNano())

	randomNumber := rand.Intn(max-min) + min

	return randomNumber
}
