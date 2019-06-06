package helpers

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
)

type Process struct {
	PID            string
	ArrivalTime    int
	ProcessingTime int
	CompletionTime int
	WaitingTime    int
}

func ReadFileProcesses(filename string) []Process {
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
			ArrivalTime:    ConvertStringToInteger(line[1]),
			ProcessingTime: ConvertStringToInteger(line[2]),
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
