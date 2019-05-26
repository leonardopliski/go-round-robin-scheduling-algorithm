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
	completionTime int
	waitingTime int
}

type executeProcessesInterface func(processes []Process)

func main() {
	processes := readFileProcesses("processes.csv")

	quantumTime := 2

	processingOrder, averageProcessingTime := executeProcessesWithRoundRobin(processes, quantumTime)

	printProcesses(processes, processingOrder, averageProcessingTime)
}

func readFileProcesses(filename string) []Process {
	reader := openFileAndGetAReader(filename)

	var processes []Process

	for {
		line, err := reader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		processes = append(processes, Process {
			PID: line[0],
			arrivalTime: convertStringToInteger(line[1]),
			processingTime: convertStringToInteger(line[2]),
		})
	}

	return processes
}

func openFileAndGetAReader(fileName string) *csv.Reader {
	csvFile, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	return csv.NewReader(bufio.NewReader(csvFile))
}

func executeProcessesWithRoundRobin (processes []Process, quantum int) (resultingProcessingOrder []string, averageProcessingTime float32) {

	currentTime := 0

	processingOrder := make([]string, 0)

	executeSchedulingWithProcessingTimeAndArrivalTimeBackup(processes, func(processes []Process) {
		for {
			completed := true
			for i := 0; i < len(processes); i++ {
				if processes[i].arrivalTime <= currentTime {
					if processes[i].arrivalTime <= quantum {
						processes[i], currentTime, completed = executeProcess(processes[i], currentTime, quantum)
						if completed == false {
							processingOrder = append(processingOrder, processes[i].PID)
						}
					} else if processes[i].arrivalTime > quantum {
						for j := 0; j < len(processes); j++ {

							if processes[j].arrivalTime < processes[i].arrivalTime {
								processes[j], currentTime, completed = executeProcess(processes[j], currentTime, quantum)
								if completed == false {
									processingOrder = append(processingOrder, processes[j].PID)
								}
							}
							processes[i], currentTime, completed = executeProcess(processes[i], currentTime, quantum)
							if completed == false {
								processingOrder = append(processingOrder, processes[i].PID)
							}
						}
					}
				} else if processes[i].arrivalTime > currentTime {
					currentTime++
					i--
				}
			}
			if completed {
				break
			}
		}
	})

	updateProcessesWaitingTime(processes)

	return processingOrder, calculateAverageProcessingTime(processes)
}

func executeSchedulingWithProcessingTimeAndArrivalTimeBackup(processes []Process, callback executeProcessesInterface) {
	processArrivalTimeBackup := make([]int, len(processes))
	processProcessingTimeBackup := make([]int, len(processes))

	for i := 0; i < len(processes); i++ {
		processArrivalTimeBackup[i] = processes[i].arrivalTime
		processProcessingTimeBackup[i] = processes[i].processingTime
	}

	callback(processes)

	for i := 0; i < len(processes); i++ {
		processes[i].arrivalTime = processArrivalTimeBackup[i]
		processes[i].processingTime = processProcessingTimeBackup[i]
	}
}

func executeProcess(process Process, currentTime int, quantum int) (updatedProcess Process, updatedTime int, processCompletedFlag bool) {
	processCompleted := Ternary(process.processingTime > 0, false, true).(bool)
	if processCompleted == false {
		if process.processingTime > quantum {
			currentTime += quantum
			process.processingTime -= quantum
			process.arrivalTime += quantum
		} else {
			currentTime += process.processingTime
			process.processingTime = 0
			process.completionTime = currentTime
		}
	}
	return process, currentTime, processCompleted
}

func printProcesses(processes []Process, finalProcessingOrder []string, totalAverageProcessingTime float32) {

	for i := 0; i < len(processes); i++ {
		fmt.Println("Waiting Time [", processes[i].PID,"]", processes[i].waitingTime, "ut")
	}

	fmt.Println("Average processing time", totalAverageProcessingTime, "ut")

	fmt.Println("Final processing order->", finalProcessingOrder)
}

func updateProcessesWaitingTime(processes []Process) {
	for i := 0; i < len(processes); i++ {
		processes[i].waitingTime = processes[i].completionTime - processes[i].arrivalTime - processes[i].processingTime
	}
}

func calculateAverageProcessingTime(processes []Process) (averageProcessingTime float32) {
	return getProcessesTotalWaitingTime(processes) / float32(len(processes))
}

func getProcessesTotalWaitingTime(processes []Process) (totalWaitingTime float32){
	var totalTime float32
	for i := 0; i < len(processes); i++ {
		totalTime += float32(processes[i].waitingTime)
	}
	return totalTime
}

func Ternary(statement bool, a, b interface{}) interface{} {
	if statement {
		return a
	}
	return b
}

func convertStringToInteger(value string) int {
	integerValue, err := strconv.Atoi(value)
	if err != nil {
		log.Fatal(err)
	}
	return integerValue
}