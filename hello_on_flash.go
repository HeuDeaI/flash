package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const allowedSerial = "9698071303404560064"

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
	switch runtime.GOOS {
	case "linux":
		out, err := exec.Command("udevadm", "info", "--query=all", "--name="+device).Output()
		if err != nil {
			return "", err
		}
		for _, line := range strings.Split(string(out), "\n") {
			if strings.Contains(line, "ID_SERIAL_SHORT=") {
				return strings.TrimSpace(strings.Split(line, "=")[1]), nil
			}
		}
	case "darwin":
		out, err := exec.Command("diskutil", "info", device).Output()
		if err != nil {
			return "", err
		}
		for _, line := range strings.Split(string(out), "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "Disk / Partition UUID:") {
				return strings.TrimSpace(strings.TrimPrefix(line, "Disk / Partition UUID:")), nil
			}
			if strings.HasPrefix(line, "Volume UUID:") {
				return strings.TrimSpace(strings.TrimPrefix(line, "Volume UUID:")), nil
			}
		}
	case "windows":
		out, err := exec.Command("wmic", "diskdrive", "get", "SerialNumber,DeviceID").Output()
		if err != nil {
			return "", err
		}
		for _, line := range strings.Split(string(out), "\n") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				return fields[len(fields)-1], nil
			}
		}
	}
	return "", fmt.Errorf("unsupported OS or device not found")
}

func main() {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		os.Exit(1)
	}
	exePath, _ = filepath.EvalSymlinks(exePath)

	device, err := getMountDevice(exePath)
	if err != nil {
		fmt.Println("Error detecting mount device:", err)
		os.Exit(1)
	}

	serial, err := getDeviceSerial(device)
	if err != nil {
		fmt.Println("Error detecting device serial:", err)
		os.Exit(1)
	}

	if serial != allowedSerial {
		fmt.Println("Unauthorized location — exiting")
		os.Exit(1)
	}

	fmt.Println("Hello World — running only from authorized flash device")
}
