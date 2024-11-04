package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

// Define a struct to hold the sensor values
type SensorData struct {
	Temperature string
	Humidity    string
	Pressure    string
	Altitude    string
}

// Global instance of the struct and a mutex for thread-safe access
var sensorData SensorData
var mu sync.Mutex

// ------------------------------------------------------------

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Initialize a slice containing the paths to the two files. Base template must be the *first* file in the slice.
	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}

	// Parse the files in the slice and store the resulting templates in a template set
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error in ParseFiles", 500)
		return
	}

	// Use the ExecuteTemplate() method to write the content of the "base" template as the response body
	mu.Lock()
	defer mu.Unlock()
	err = ts.ExecuteTemplate(w, "base", sensorData)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error in ExecuteTemplate", 500)
	}
}

// ------------------------------------------------------------

func createSensorData(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed"))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest) // Response status: 400 Bad Request
		return
	}
	defer r.Body.Close()

	// Convert the body to a string
	bodyString := string(body)
	log.Println("Received data:", bodyString)

	// Split the string by commas to separate the sensor values
	sensorValues := strings.Split(bodyString, ",")
	log.Println("Parsed sensor values:", sensorValues)
	if len(sensorValues) < 4 {
		http.Error(w, "Invalid data format", http.StatusBadRequest) // Response body: Invalid data format
		return
	}

	// Extract the sensor values
	mu.Lock()
	sensorData = SensorData{
		Temperature: strings.TrimSpace(sensorValues[0]),
		Humidity:    strings.TrimSpace(sensorValues[1]),
		Pressure:    strings.TrimSpace(sensorValues[2]),
		Altitude:    strings.TrimSpace(sensorValues[3]),
	}
	mu.Unlock()

	log.Println("Sensor values:", sensorData)

	// Send a response to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Sensor data received successfully"))
}

// ------------------------------------------------------------

func getSensorData(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sensorData)
}

func main() {
	// Create a new ServeMux and register the routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/api/v1/sensors", createSensorData)
	mux.HandleFunc("/api/v1/sensors/data", getSensorData)

	// Start the server
	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}

