package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
)

const (
	numThreads = 60
	bufferSize = 20 * 1024 * 1024 // 20MB
)

type logEntry struct {
	line string
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	startTimeStr := flag.String("start", "", "Start time in format MM/DD/YYYY HH:MM:SS")
	endTimeStr := flag.String("end", "", "End time in format MM/DD/YYYY HH:MM:SS")
	logFilePath := flag.String("file", "", "Path to the log file")
	flag.Parse()

	if *startTimeStr == "" || *endTimeStr == "" || *logFilePath == "" {
		fmt.Println("Usage: ./cut -start=\"MM/DD/YYYY HH:MM:SS\" -end=\"MM/DD/YYYY HH:MM:SS\" -file=\"path/to/logfile.log\"")
		return
	}

	startTime := *startTimeStr
	endTime := *endTimeStr
	file, err := os.Open(*logFilePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return
	}
	fileSize := fileInfo.Size()

	chunkSize := fileSize / numThreads

	linesChannel := make(chan logEntry, 100)
	var wg sync.WaitGroup

	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			processChunk(*logFilePath, id, chunkSize, startTime, endTime, linesChannel)
		}(i)
	}

	go func() {
		wg.Wait()
		close(linesChannel)
	}()

	for entry := range linesChannel {
		fmt.Print(entry.line)
	}
}

func processChunk(filePath string, id int, chunkSize int64, startTime, endTime string, linesChannel chan logEntry) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	start := int64(id) * chunkSize
	end := start + chunkSize
	file.Seek(start, 0)
	reader := bufio.NewReaderSize(file, bufferSize)

	var buffer strings.Builder
	var capture bool

	for {
		pos, _ := file.Seek(0, 1)
		if pos >= end {
			break
		}

		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		if strings.Contains(line, startTime) {
			capture = true
		}
		if capture {
			buffer.WriteString(line)
		}
		if strings.Contains(line, endTime) {
			capture = false
			break
		}
	}

	if buffer.Len() > 0 {
		linesChannel <- logEntry{line: buffer.String()}
	}
}
