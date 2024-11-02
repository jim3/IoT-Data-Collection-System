package main

import (
	"log"
	"time"

	"go.bug.st/serial"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
)

type myService struct{}

// Execute implements the svc.Handler interface
func (m *myService) Execute(args []string, r <-chan svc.ChangeRequest, s chan<- svc.Status) (bool, uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	s <- svc.Status{State: svc.StartPending}
	go m.run()
	s <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Stop, svc.Shutdown:
				break loop
			default:
				continue loop
			}
		}
	}
	s <- svc.Status{State: svc.StopPending}
	return false, 0
}

func (m *myService) run() {
	for {
		GetPortsList()
		mode := &serial.Mode{
			BaudRate: 115200,
		}
		SerialPortConfig("COM3", mode)
		time.Sleep(10 * time.Second)
	}
}

func main() {
	isInteractive, err := svc.IsAnInteractiveSession()
	if err != nil {
		log.Fatalf("failed to determine if we are running in an interactive session: %v", err)
	}
	if !isInteractive {
		runService("GoSerialService", false)
		return
	}
}

func runService(name string, isDebug bool) {
	const svcName = "GoSerialService"
	if isDebug {
		log.Printf("Starting %s service in debug mode.", svcName)
	} else {
		log.Printf("Starting %s service.", svcName)
	}
	run := svc.Run
	if isDebug {
		run = debug.Run
	}
	err := run(svcName, &myService{})
	if err != nil {
		log.Fatalf("%s service failed: %v", svcName, err)
	}
	log.Printf("%s service stopped.", svcName)
}
