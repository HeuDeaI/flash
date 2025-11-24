package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func getSerial(device string) (string, error) {

	if runtime.GOOS != "linux" {
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	out, err := exec.Command("udevadm", "info", "--query=all", "--name="+device).Output()
	if err != nil {
		return "", fmt.Errorf("no device found at %s", device)
	}

	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "ID_SERIAL_SHORT=") {
			return strings.TrimSpace(strings.Split(line, "=")[1]), nil
		}
	}

	return "", fmt.Errorf("no device serial found for %s", device)
}

func main() {

	allowed := []string{"9698071303404560064", "1234567890ABCDEF"}

	devices := []string{"/dev/sdb", "/dev/sdc", "/dev/sdd"}

	authorized := false

	for _, device := range devices {
		serial, err := getSerial(device)
		if err != nil {
			fmt.Println("Info:", err)
			continue
		}

		fmt.Println("Device:", device, "Serial:", serial)

		for _, a := range allowed {
			if serial == a {
				authorized = true
				break
			}
		}
	}

	if !authorized {
		fmt.Println("Unauthorized device: no matching ID")
		return
	}

	fmt.Println("Hello World")
}
