// package main

// import (
// 	"fmt"
// 	"os"
// 	"time"
// )

// func check(e error) {
// 	if e != nil {
// 		panic(e)
// 	}
// }

// func WatchDog(filename string, sig chan string) {

// 	interval := 5 * time.Second

// 	for {
// 		out, err := os.ReadFile(filename)
// 		check(err)
// 		time.Sleep(interval)
// 		sig <- string(out[:])
// 	}
// }

// func RunNextflow(workflow) {

// }

// func main() {

// 	messages := make(chan string)

// 	go WatchDog("something.txt", messages)

// 	for n := range messages {
// 		fmt.Println(n)
// 	}
// }

package main

import (
	"fmt"
	"math/rand"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// The type for channel communication
type Work = string
type PID = int

func readCommands(commandChan chan<- Work) {

	rand.Seed(time.Now().UnixNano())
	min := 5
	max := 10
	var command Work

	for {
		// Read commands from file

		integ := rand.Intn(max-min+1) + min

		command = fmt.Sprintf("sleep %v", integ)

		// Add the command to the channel
		commandChan <- command

		time.Sleep(10 * time.Second)
	}
}

func runCommand(commandChan <-chan Work, processChan chan PID) {

	for command := range commandChan {

		// Run the command
		fmt.Println("Starting command:", command)
		cmd := exec.Command("bash", "-c", command)
		err := cmd.Start()
		if err != nil {
			fmt.Println("Error starting command:", err)
			continue
		}

		// Add the PID to the channel
		processChan <- cmd.Process.Pid

	}

	// cmd := exec.Command("nextflow", "run", pipeline)
	// err := cmd.Start()
	// if err != nil {
	// 	return -1, fmt.Errorf("Error starting Nextflow pipeline: %v", err)
	// }
	// return cmd.Process.Pid, nil

}

func monitorProcess(processChan <-chan PID, reportChan chan<- PID) {

	PIDs := make(map[int]int)

	go func() {
		for {
			pid := <-processChan
			PIDs[pid] = 1
		}
	}()

	for range time.Tick(10 * time.Second) {
		for pid, _ := range PIDs {
			output, err := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "state").Output()
			if err != nil {
				fmt.Println("Error getting process status:", err)
				continue
			}

			status := strings.TrimSpace(string(output))
			if status == "Z" {
				fmt.Println("Nextflow pipeline failed")
				continue
			} else if status == "R" {
				fmt.Println("Nextflow pipeline running")
			} else {
				fmt.Println("Nextflow pipeline completed successfully")
				// Remove PID from PIDs
				delete(PIDs, pid)
				continue
			}
		}
	}
}

func main() {

	commandChan := make(chan string)
	PIDChan := make(chan int)
	reportChan := make(chan int)

	go readCommands(commandChan)
	go runCommand(commandChan, PIDChan)
	go monitorProcess(PIDChan, reportChan)

	for {
		<-reportChan

	}

	// What am I supposed to do here?
}
