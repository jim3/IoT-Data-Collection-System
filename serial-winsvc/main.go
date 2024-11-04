package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"go.bug.st/serial"
)

const url string = "https://jim3.xyz/api/v1/sensors"

func GetPortsList() {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}
	for _, port := range ports {
		fmt.Printf("Found port: %v\n", port)
	}
}

func SerialPortConfig(portName string, mode *serial.Mode) {
	port, err := serial.Open(portName, mode)
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()

	buff := make([]byte, 100)
	var dataBuffer bytes.Buffer

	for {
		n, err := port.Read(buff)
		if err != nil {
			log.Fatal(err)
			break
		}
		if n == 0 {
			fmt.Println("\nEOF")
			break
		}

		dataBuffer.Write(buff[:n])

		data := dataBuffer.String()
		if data[len(data)-1] == '\n' {
			fmt.Println("Complete set of sensor values =>", data)
			HttpPostRequest(data)
			dataBuffer.Reset()
		}
	}
}

func HttpPostRequest(data string) {
	resp, err := http.Post(url, "application/octet-stream", bytes.NewBufferString(data))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Response body:", string(body))
}
