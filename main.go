package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type Process struct {
	PID string
	arrivalTime  int
	processingTime int
	priority int
	completionTime int
	waitingTime int
}

func main() {
	csvFile, _ := os.Open("processes.csv")

	reader := csv.NewReader(bufio.NewReader(csvFile))
	var processes []Process

	for {
		line, e := reader.Read()

		if e == io.EOF {
			break
		} else if e != nil {
			log.Fatal(e)
		}

		arrivalTime, e := strconv.Atoi(line[0])
		if e != nil {
			log.Fatal(e)
		}

		priority, e := strconv.Atoi(line[2])
		if e != nil {
			log.Fatal(e)
		}

		processingTime, e := strconv.Atoi(line[3])
		if e != nil {
			log.Fatal(e)
		}

		processes = append(processes, Process {
			PID: line[1],
			arrivalTime: arrivalTime,
			processingTime: processingTime,
			priority: priority,
		})
	}

	quantum := 2

	calculateRoundRobin(processes, quantum)
}

func calculateRoundRobin (processes []Process, quantum int) {

	for i := 0; i < calculateMaximumProcessingTime(processes); i++ {

		for j := i; j < len(processes); j++ {
			if processes[j].arrivalTime <= i {

				if processes[j].processingTime > 0 {
					if processes[j].processingTime <= quantum {
						processes[j].completionTime = i + processes[j].processingTime
						processes[j].processingTime = 0
					} else {
						processes[j].processingTime -= quantum
					}
				}

			}
		}

	}

	for i := 0; i < len(processes); i++ {
		fmt.Println(processes[i])
	}
	//for i := 0; i <= maximumProcessingTime; i++ {
	//	for j := 0; j < len(processes); j++ {
	//		if processes[j].arrivalTime <= i {
	//
	//			if processes[j].processingTime > quantum {
	//				processes[j].processingTime -= quantum
	//				processes[j].arrivalTime += quantum
	//			} else if processes[j].processingTime > 0 {
	//				processes[j].arrivalTime += processes[j].processingTime
	//				processes[j].processingTime = 0
	//				processes[j].completionTime = processes[j].arrivalTime
	//			}
	//
	//		}
	//	}
	//}
	//fmt.Println(processes)
}

func calculateMaximumProcessingTime(processes []Process) int {
	var total int = 0
	for _, process := range processes {
		total += process.processingTime
	}
	return total
}

func calculateAverageProcessingTime(processingTimes []int) int {
	var sum int = 0
	for _, i := range processingTimes {
		sum += i
	}
	return (sum / len(processingTimes))
}
func calculateProcessesWaitingTime(processes []string, processingTimes []int, arrivalTime []int, quantum int) []int {
	remainingProcessingTime := make([]int, len(processes))

	copy(remainingProcessingTime, processingTimes)

	proccessWaitingTime := make([]int, len(processes))
	proccessArrivalTime := make([]int, len(processes))

	currentUnitTime := 0

	for {
		completed := true

		for i := 0; i < len(processes); i++ {
			if remainingProcessingTime[i] > 0 {
				completed = false
				if remainingProcessingTime[i] > quantum {
					//currentUnitTime = quantum

					currentUnitTime += quantum
					remainingProcessingTime[i] -= quantum

				} else {
					currentUnitTime += remainingProcessingTime[i]

					proccessWaitingTime[i] = currentUnitTime - processingTimes[i]

					remainingProcessingTime[i] = 0

					proccessArrivalTime[i] = currentUnitTime
				}
			}
		}

		if completed == true {
			break
		}
	}

	return proccessWaitingTime
}

func calculateTurnAroundTime(processingTimes []int, waitingTimes []int) []int {
	i := 0

	turnAroundTime := make([]int, len(processingTimes))

	for i = 0; i < len(processingTimes); i++ {
		turnAroundTime[i] = processingTimes[i] + waitingTimes[i]
	}

	return turnAroundTime
}

func valueInArray(a int, array []int) bool {
	for _, value := range array {
		if value == a {
			return true
		}
	}
	return false
}