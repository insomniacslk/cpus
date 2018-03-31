package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

// TODO daemon mode, that randomizes periodically the CPUs that are on and off

var _availableCPUs []int

func getCPUs() ([]int, error) {
	if _availableCPUs != nil {
		return _availableCPUs, nil
	}
	files, err := ioutil.ReadDir("/sys/devices/system/cpu")
	if err != nil {
		return nil, err
	}
	cpus := make([]int, 0)
	for _, f := range files {
		if !strings.HasPrefix(f.Name(), "cpu") {
			continue
		}
		cpunum, err := strconv.ParseInt(f.Name()[3:], 10, 16)
		if err != nil {
			continue
		}
		cpus = append(cpus, int(cpunum))
	}
	_availableCPUs = cpus
	return _availableCPUs, nil
}

func getCPUStatus(cpunum int) (bool, error) {
	if cpunum == 0 {
		// CPU0 is always online
		return true, nil
	}
	data, err := ioutil.ReadFile(fmt.Sprintf("/sys/devices/system/cpu/cpu%d/online", cpunum))
	if err != nil {
		return false, err
	}
	online, err := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 8)
	if err != nil {
		return false, err
	}
	if online == 0 {
		return false, nil
	}
	return true, nil
}

func showCPUs() {
	cpus, err := getCPUs()
	if err != nil {
		log.Printf("Failed to get CPUs: %v", err)
		return
	}
	for _, cpunum := range cpus {
		online, err := getCPUStatus(cpunum)
		if err != nil {
			log.Printf("Failed to get CPU status for CPU%d: %v", cpunum, err)
			return
		}
		var status string
		if online {
			status = "online"
		} else {
			status = "offline"
		}
		fmt.Printf("CPU%d is %s\n", cpunum, status)
	}
}

func IsValidCPU(cpunum int) bool {
	availableCPUs, err := getCPUs()
	if err != nil {
		// I don't like using panic in a library function but..
		panic(err)
	}
	for _, avCpu := range availableCPUs {
		if cpunum == avCpu {
			return true
		}
	}
	return false
}

func changeCPUStatus(status bool, cpus ...int) error {
	var (
		statusStr string
		online    []byte
	)
	if status {
		online = []byte("1")
		statusStr = "online"
	} else {
		online = []byte("0")
		statusStr = "offline"
	}
	fmt.Printf("Changing status for cpus: %v to %s\n", cpus, statusStr)
	for _, cpu := range cpus {
		if !IsValidCPU(cpu) {
			return fmt.Errorf("Invalid CPU number: %d", cpu)
		}
		if cpu == 0 {
			// skip CPU0, it's always online
			continue
		}
		err := ioutil.WriteFile(
			fmt.Sprintf("/sys/devices/system/cpu/cpu%d/online", cpu),
			online,
			0644,
		)
		if err != nil {
			return fmt.Errorf("Cannot change CPU %d status to %s: %v", cpu, statusStr, err)
		}
	}
	return nil
}

func cmdCPUOn(cpus ...int) error {
	return changeCPUStatus(true, cpus...)
}

func cmdCPUOff(cpus ...int) error {
	return changeCPUStatus(false, cpus...)
}

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Println("    status : print CPU online status. If no argument is passed, this is the default")
		fmt.Println("    on     : turn CPUs on. Optionally pass a list of CPU # to turn on selectively")
		fmt.Println("    off    : turn CPUs off. Optionally pass a list of CPU # to turn off selectively")
	}
	flag.Parse()
	var (
		cmd string
		err error
	)
	if len(flag.Args()) == 0 {
		cmd = "status"
	} else {
		cmd = flag.Arg(0)
	}
	if cmd == "status" {
		showCPUs()
	} else if cmd == "on" || cmd == "off" {
		cpus := make([]int, 0)
		if len(flag.Args()[1:]) == 0 {
			// set status for all CPUs
			cpus, err = getCPUs()
			if err != nil {
				log.Fatalf("Cannot get CPUs: %v", err)
			}
		} else {
			for _, cpuStr := range flag.Args()[1:] {
				cpunum, err := strconv.ParseInt(cpuStr, 10, 16)
				if err != nil {
					log.Fatalf("Invalid CPU number: %v", cpuStr)
				}
				cpus = append(cpus, int(cpunum))
			}
		}
		var err error
		if cmd == "on" {
			err = cmdCPUOn(cpus...)
		} else {
			err = cmdCPUOff(cpus...)
		}
		if err != nil {
			log.Fatalf("Cannot change CPU status: %v", err)
		}
		showCPUs()
	} else {
		log.Fatalf("Unknown command: %s", cmd)
	}
}
