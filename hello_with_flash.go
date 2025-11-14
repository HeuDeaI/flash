package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func getSerial() (string, error) {
	switch runtime.GOOS {
	case "linux":
		out, err := exec.Command("udevadm", "info", "--query=all", "--name=/dev/sdc").Output()
		if err != nil {
			return "", err
		}
		for _, line := range strings.Split(string(out), "\n") {
			if strings.Contains(line, "ID_SERIAL_SHORT=") {
				return strings.TrimSpace(strings.Split(line, "=")[1]), nil
			}
		}
	case "darwin":
		out, err := exec.Command("diskutil", "info", "/dev/disk2").Output()
		if err != nil {
			return "", err
		}
		return parseMacSerial(string(out)), nil
	case "windows":
		out, err := exec.Command("wmic", "diskdrive", "get", "SerialNumber").Output()
		if err != nil {
			return "", err
		}
		return parseWinSerial(string(out)), nil
	}
	return "", fmt.Errorf("unsupported OS")
}

func main() {
	allowed := "9698071303404560064"
	serial, err := getSerial()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if serial != allowed {
		fmt.Println("Unauthorized device")
		return
	}
	fmt.Println("Hello World")
}

func parseMacSerial(output string) string {
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Disk / Partition UUID:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "Disk / Partition UUID:"))
		}
		if strings.HasPrefix(line, "Volume UUID:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "Volume UUID:"))
		}
	}
	return ""
}

func parseWinSerial(output string) string {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			return fields[len(fields)-1]
		}
	}
	return ""
}
