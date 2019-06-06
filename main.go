package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"./helpers"
)

type Process struct {
	PID            string
	arrivalTime    int
	processingTime int
	completionTime int
	waitingTime    int
}

type executeProcessesInterface func(processes []Process)
type processLoopInnerFunctionInterface func(process *Process) (walkBackToPriorProcess bool)

func main() {
	processes := readFileProcesses("processes.csv")

	quantumTime := 2

	processingOrder, averageProcessingTime := executeProcessesWithRoundRobinTimeScheduling(processes, quantumTime)

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

		processes = append(processes, Process{
			PID:            line[0],
			arrivalTime:    helpers.ConvertStringToInteger(line[1]),
			processingTime: helpers.ConvertStringToInteger(line[2]),
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

func executeProcessesWithRoundRobinTimeScheduling(processes []Process, quantum int) (resultingProcessingOrder []string, averageProcessingTime float32) {

	processingOrder := make([]string, 0)

	executeSchedulingWithProcessingTimeAndArrivalTimeBackup(processes, func(processes []Process) {
		currentTime := 0

		for {
			completed := true
			loopThroughProcesses(processes, func(currentProcess *Process) (walkBackToPriorProcessFlag bool) {
				if currentProcess.arrivalTime <= currentTime {
					if currentProcess.arrivalTime <= quantum {
						currentProcess, currentTime, completed = executeProcess(currentProcess, currentTime, quantum)
						if completed == false {
							processingOrder = appendProcessToProcessingOrder(currentProcess, processingOrder)
						}
					} else if currentProcess.arrivalTime > quantum {
						loopThroughProcesses(processes, func(process *Process) (walkBackToPriorProcessFlag bool) {
							if process.arrivalTime < currentProcess.arrivalTime {
								process, currentTime, completed = executeProcess(process, currentTime, quantum)
								if completed == false {
									processingOrder = appendProcessToProcessingOrder(process, processingOrder)
								}
							}
							currentProcess, currentTime, completed = executeProcess(currentProcess, currentTime, quantum)
							if completed == false {
								processingOrder = appendProcessToProcessingOrder(currentProcess, processingOrder)
							}
							return false
						})
					}
				} else if currentProcess.arrivalTime > currentTime {
					currentTime++
					return true
				}
				return false
			})
			if completed {
				break
			}
		}
	})

	updateProcessesWaitingTime(processes)

	return processingOrder, calculateAverageProcessingTime(processes)
}

func appendProcessToProcessingOrder(process *Process, processingOrder []string) []string {
	return append(processingOrder, process.PID)
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

func loopThroughProcesses(processes []Process, callback processLoopInnerFunctionInterface) {
	for i := 0; i < len(processes); i++ {
		walkToPriorProcessFlag := callback(&processes[i])
		if walkToPriorProcessFlag {
			i--
		}
	}
}

func executeProcess(process *Process, currentTime int, quantum int) (updatedProcess *Process, updatedTime int, processCompletedFlag bool) {
	processCompleted := helpers.Ternary(process.processingTime > 0, false, true).(bool)
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
		fmt.Println("Waiting Time [", processes[i].PID, "]", processes[i].waitingTime, "ut")
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

func getProcessesTotalWaitingTime(processes []Process) (totalWaitingTime float32) {
	var totalTime float32
	for i := 0; i < len(processes); i++ {
		totalTime += float32(processes[i].waitingTime)
	}
	return totalTime
}