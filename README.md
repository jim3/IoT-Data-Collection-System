### IoT Data Collection System

Many different components that work together to collect sensor data from a Waveshare BME680 sensor and display in real-time. The components are:

- An `ESP32 Mini-1` microcontroller that uses [MicroPython](https://micropython.org) to read sensor data from a **Waveshare BME680** environmental sensor. 
The sensor measures temperature, humidity, pressure, and gas resistance. The sensor is connected to the ESP32 via I2C. The ESP32 reads the sensor data and sends it out via a serial port.
All of the hard work is done with a MicroPython driver for the BME680 sensor located here [Micropython Driver for a BME680 breakout](https://github.com/adafruit/Adafruit_BME680) 
and here [Micropython Driver](https://github.com/robert-hh/BME680-Micropython) On the ESP32 Mini-1, you can use PIN's 21 and 22 for the sensor.
  
- Code that reads the sensor data from a Windows COM port (the ESP32 is connected to the computer via USB) and sends it to a live server via a `POST` request. The script uses the wonderful [serial](https://github.com/bugst/go-serial) package to read from the serial port. The [serial package](https://pkg.go.dev/go.bug.st/serial) documentation is very helpful and the package is *very easy* to use. One other thing to note is I wanted to do this in Go *and* use Windows. One huge obstacle I run into was trying to figure out how to run this code locally and as a Windows service. I found an EXTREMELY useful package provided by Go called https://pkg.go.dev/golang.org/x/sys/windows/svc. This package allows you to run Go code as a Windows service! 😲 I was able to run this code as a service and it worked perfectly! A lot of help was provided for this step but I finally got it working! 😎

- Server-side [Go](https://go.dev) uses Linux. It listens for the incoming connections (sensor data), processes it, and sends it to a client to display the IoT data in real-time.

A demo of this lives at https://jim3.xyz

---

### Usage

1. Connect the ESP32 to the BME680 sensor via I2C and connect the ESP32 to your computer via USB.

2. Start the server: `go run main.go`

3. Start the client script: `go run main.go`

4. Open the client in your browser: `http://localhost:8080`


---

### To Do

I'll try and make the README.md more detailed and add more information about the project but after teh amount of time this took I'm just happy it works. 😅😅😅
