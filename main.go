package main

import (
	"./helpers"
	"./algorithms"
)

func main() {
	processes := helpers.ReadFileProcesses("processes.csv")

	quantumTime := 2

	processingOrder, averageProcessingTime := algorithms.ExecuteProcessesWithRoundRobinTimeScheduling(processes, quantumTime)

	algorithms.PrintProcesses(processes, processingOrder, averageProcessingTime)
}