package algorithms

import (
	"fmt"
	"../helpers"
)

type Process struct {
	PID            string
	arrivalTime    int
	processingTime int
	completionTime int
	waitingTime    int
}

type executeProcessesInterface func(processes []helpers.Process)
type processLoopInnerFunctionInterface func(process *helpers.Process) (walkBackToPriorProcess bool)

func ExecuteProcessesWithRoundRobinTimeScheduling(processes []helpers.Process, quantum int) (resultingProcessingOrder []string, averageProcessingTime float32) {

	processingOrder := make([]string, 0)

	executeSchedulingWithProcessingTimeAndArrivalTimeBackup(processes, func(processes []helpers.Process) {
		currentTime := 0

		for {
			completed := true
			loopThroughProcesses(processes, func(currentProcess *helpers.Process) (walkBackToPriorProcessFlag bool) {
				if currentProcess.ArrivalTime <= currentTime {
					if currentProcess.ArrivalTime <= quantum {
						currentProcess, currentTime, completed = executeProcess(currentProcess, currentTime, quantum)
						if completed == false {
							processingOrder = appendProcessToProcessingOrder(currentProcess, processingOrder)
						}
					} else if currentProcess.ArrivalTime > quantum {
						loopThroughProcesses(processes, func(process *helpers.Process) (walkBackToPriorProcessFlag bool) {
							if process.ArrivalTime < currentProcess.ArrivalTime {
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
				} else if currentProcess.ArrivalTime > currentTime {
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

func appendProcessToProcessingOrder(process *helpers.Process, processingOrder []string) []string {
	return append(processingOrder, process.PID)
}

func executeSchedulingWithProcessingTimeAndArrivalTimeBackup(processes []helpers.Process, callback executeProcessesInterface) {
	processArrivalTimeBackup := make([]int, len(processes))
	processProcessingTimeBackup := make([]int, len(processes))

	for i := 0; i < len(processes); i++ {
		processArrivalTimeBackup[i] = processes[i].ArrivalTime
		processProcessingTimeBackup[i] = processes[i].ProcessingTime
	}

	callback(processes)

	for i := 0; i < len(processes); i++ {
		processes[i].ArrivalTime = processArrivalTimeBackup[i]
		processes[i].ProcessingTime = processProcessingTimeBackup[i]
	}
}

func loopThroughProcesses(processes []helpers.Process, callback processLoopInnerFunctionInterface) {
	for i := 0; i < len(processes); i++ {
		walkToPriorProcessFlag := callback(&processes[i])
		if walkToPriorProcessFlag {
			i--
		}
	}
}

func executeProcess(process *helpers.Process, currentTime int, quantum int) (updatedProcess *helpers.Process, updatedTime int, processCompletedFlag bool) {
	processCompleted := helpers.Ternary(process.ProcessingTime > 0, false, true).(bool)
	if processCompleted == false {
		if process.ProcessingTime > quantum {
			currentTime += quantum
			process.ProcessingTime -= quantum
			process.ArrivalTime += quantum
		} else {
			currentTime += process.ProcessingTime
			process.ProcessingTime = 0
			process.CompletionTime = currentTime
		}
	}
	return process, currentTime, processCompleted
}

func PrintProcesses(processes []helpers.Process, finalProcessingOrder []string, totalAverageProcessingTime float32) {

	for i := 0; i < len(processes); i++ {
		fmt.Println("Waiting Time [", processes[i].PID, "]", processes[i].WaitingTime, "ut")
	}

	fmt.Println("Average processing time", totalAverageProcessingTime, "ut")

	fmt.Println("Final processing order->", finalProcessingOrder)
}

func updateProcessesWaitingTime(processes []helpers.Process) {
	for i := 0; i < len(processes); i++ {
		processes[i].WaitingTime = processes[i].CompletionTime - processes[i].ArrivalTime - processes[i].ProcessingTime
	}
}

func calculateAverageProcessingTime(processes []helpers.Process) (averageProcessingTime float32) {
	return getProcessesTotalWaitingTime(processes) / float32(len(processes))
}

func getProcessesTotalWaitingTime(processes []helpers.Process) (totalWaitingTime float32) {
	var totalTime float32
	for i := 0; i < len(processes); i++ {
		totalTime += float32(processes[i].WaitingTime)
	}
	return totalTime
}