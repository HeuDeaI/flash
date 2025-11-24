package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// getMountDevice returns the block device backing a given path
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

// normalize partition device (/dev/sda2 -> /dev/sda, /dev/nvme0n1p3 -> /dev/nvme0n1)
func getParentDevice(device string) string {
	// regex handles both sdXn and nvmeXpYn style names
	re := regexp.MustCompile(`^(/dev/[a-zA-Z]+[0-9]+)`)
	// strip trailing digits/pN
	if strings.HasPrefix(device, "/dev/nvme") {
		// remove partition suffix like p1, p2
		re = regexp.MustCompile(`^(/dev/nvme[0-9]+n[0-9]+)p?[0-9]*$`)
	} else {
		re = regexp.MustCompile(`^(/dev/[a-z]+[a-z][0-9]*)[0-9]*$`)
	}
	m := re.FindStringSubmatch(device)
	if len(m) > 1 {
		return m[1]
	}
	return device
}

// query udevadm for serial
func getDeviceSerial(device string) (string, error) {
	out, err := exec.Command("udevadm", "info", "--query=all", "--name="+device).Output()
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "ID_SERIAL_SHORT=") {
			return strings.TrimSpace(strings.Split(line, "=")[1]), nil
		}
	}
	return "", fmt.Errorf("serial not found for %s", device)
}

func main() {
	allowed := "9698071303404560064"

	// Example: check the device mounted at /home
	device, err := getMountDevice("/home")
	if err != nil {
		fmt.Println("Error getting mount device:", err)
		return
	}

	parent := getParentDevice(device)
	serial, err := getDeviceSerial(parent)
	if err != nil {
		fmt.Println("Error getting device serial:", err)
		return
	}

	fmt.Printf("Device: %s (parent: %s)\n", device, parent)
	fmt.Printf("Serial: %s\n", serial)

	if serial != allowed {
		fmt.Println("Unauthorized device")
		return
	}
	fmt.Println("Hello World")
}
