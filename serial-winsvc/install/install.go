package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/svc/mgr"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: installer.exe <install|remove>")
		return
	}

	// Get the path to the service executable
	servicePath := filepath.Join(filepath.Dir(os.Args[0]), "GoSerialService.exe")
	if _, err := os.Stat(servicePath); os.IsNotExist(err) {
		log.Fatalf("Service executable not found at: %s", servicePath)
	}

	cmd := os.Args[1]
	switch cmd {
	case "install":
		err := installService("GoSerialService", "Go Serial Service", servicePath)
		if err != nil {
			log.Fatalf("Failed to install service: %v", err)
		}
		fmt.Println("Service installed successfully")
	case "remove":
		err := removeService("GoSerialService")
		if err != nil {
			log.Fatalf("Failed to remove service: %v", err)
		}
		fmt.Println("Service removed successfully")
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		fmt.Println("Usage: installer.exe <install|remove>")
	}
}

func installService(name, desc, exepath string) error {
	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to service manager: %v", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(name)
	if err == nil {
		s.Close()
		return fmt.Errorf("service %s already exists", name)
	}

	s, err = m.CreateService(name, exepath, mgr.Config{
		DisplayName: desc,
		StartType:   mgr.StartAutomatic,
	})
	if err != nil {
		return fmt.Errorf("failed to create service: %v", err)
	}
	defer s.Close()
	return nil
}

func removeService(name string) error {
	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to service manager: %v", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("service %s not found", name)
	}
	defer s.Close()

	err = s.Delete()
	if err != nil {
		return fmt.Errorf("failed to delete service: %v", err)
	}
	return nil
}
