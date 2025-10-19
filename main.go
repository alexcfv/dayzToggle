package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	evdev "github.com/gvalkov/golang-evdev"
)

const (
	iface     = "wlan0"
	rate      = "500kbit"
	delaySpec = "300ms"
	jitter    = "100ms"
	loss      = "3%"
	kbdDev    = "/dev/input/event11" // keyboard
)

func runCmd(cmd string, args ...string) {
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	_ = c.Run()
}

func enableLag() {
	runCmd("sudo", "tc", "qdisc", "del", "dev", iface, "root")
	runCmd("sudo", "tc", "qdisc", "add", "dev", iface, "root", "handle", "1:", "htb", "default", "1")
	runCmd("sudo", "tc", "class", "add", "dev", iface, "parent", "1:", "classid", "1:1", "htb", "rate", rate)
	runCmd("sudo", "tc", "qdisc", "add", "dev", iface, "parent", "1:1", "handle", "10:", "netem",
		"delay", delaySpec, jitter, "loss", loss)
	fmt.Println("🐢 Lag enable")
}

func disableLag() {
	runCmd("sudo", "tc", "qdisc", "del", "dev", iface, "root")
	fmt.Println("⚡ Lag disable")
}

func main() {
	dev, err := evdev.Open(kbdDev)
	if err != nil {
		log.Fatalf("Error to open keyboard: %v", err)
	}

	fmt.Println("X for lag, Z exit")

	lagActive := false

	for {
		events, err := dev.Read()
		if err != nil {
			time.Sleep(10 * time.Millisecond)
			continue
		}

		for _, e := range events {
			if e.Type == evdev.EV_KEY {
				switch e.Code {
				case evdev.KEY_X:
					if e.Value == 1 && !lagActive {
						enableLag()
						lagActive = true
					}
					if e.Value == 0 && lagActive { 
						disableLag()
						lagActive = false
					}
				case evdev.KEY_Z:
					if e.Value == 1 {
						if lagActive {
							disableLag()
						}
						fmt.Println("Exit...")
						return
					}
				}
			}
		}

		time.Sleep(10 * time.Millisecond)
	}
}

