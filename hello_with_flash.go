package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func getMountDevice(path string) (string, error) {
	out, err := exec.Command("df", "--output=source", path).Output()
	if err != nil {
		return "", err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("unexpected df output")
	}
	device := strings.TrimSpace(lines[1])
	return device, nil
}

func getDeviceSerial(device string) (string, error) {
	// Linux-only implementation
	out, err := exec.Command("udevadm", "info", "--query=all", "--name="+device).Output()
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "ID_SERIAL_SHORT=") {
			return strings.TrimSpace(strings.Split(line, "=")[1]), nil
		}
	}
	return "", fmt.Errorf("serial not found")
}

func main() {
	allowed := "9698071303404560064"

	// Example: check the device mounted at /home
	device, err := getMountDevice("/home")
	if err != nil {
		fmt.Println("Error getting mount device:", err)
		return
	}

	serial, err := getDeviceSerial(device)
	if err != nil {
		fmt.Println("Error getting device serial:", err)
		return
	}

	if serial != allowed {
		fmt.Println("Unauthorized device")
		return
	}
	fmt.Println("Hello World")
}
