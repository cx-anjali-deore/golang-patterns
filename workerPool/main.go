package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

func main() {
	const message = 100
	const workers = 3
	Dispatcher(message, workers)
}

type Message struct {
	Id      int
	Message string
}

type Result struct {
	Id               int
	ProcessedMessage string
	timeStamp        time.Time
}

func workers(workerWg *sync.WaitGroup, messageChan <-chan Message, resultsChan chan<- Result) {
	defer workerWg.Done()
	for message := range messageChan {
		time.Sleep(100 * time.Millisecond)
		resultsChan <- Result{message.Id, message.Message, time.Now()}
	}
}

func collectingResults(resultsChan <-chan Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for result := range resultsChan {
		fmt.Printf("Received message: %s  At time %s\n", result.ProcessedMessage, result.timeStamp)
	}
}

func Dispatcher(message int, worker int) {

	messagesChan := make(chan Message, message)
	resultsChan := make(chan Result, message)

	// generating the Messages in Message channel

	for i := 1; i <= message; i++ {
		messagesChan <- Message{
			Id:      i,
			Message: "Hello " + strconv.Itoa(i),
		}
	}
	var workerWg sync.WaitGroup
	workerWg.Add(worker)
	for j := 0; j < worker; j++ {
		go workers(&workerWg, messagesChan, resultsChan)
	}

	var resultWg sync.WaitGroup

	resultWg.Add(1)
	go collectingResults(resultsChan, &resultWg)

	close(messagesChan)
	workerWg.Wait()

	close(resultsChan)

	resultWg.Wait()

}
