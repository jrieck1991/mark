package main

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// n workers each in their own goroutine
// reading n messages, each message must be deleted within n time or its repeated.

const (

	// writer config
	responseBatchSize int = 10
	msgBatchSize      int = 10

	// metrics tags
	namespace string = "app"
	writer    string = "writer"
	reader    string = "reader"
)

type Response struct {
	Messages []Message
}

type Message struct {
	Data interface{}
}

func main() {
	fmt.Println("start")

	var wg sync.WaitGroup

	// start metrics http server
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":9000", nil)
	}(&wg)

	// init data stream
	dataStream := make(chan []Response, 100)

	// write stream
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		writeStream(dataStream)
	}(&wg)

	// read from stream
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		readStream(dataStream)
	}(&wg)

	fmt.Println("init complete")
	wg.Wait()
}

// deleteMessage sleeps for 300ms, a network api call
func deleteMessage(msg Message) {
	time.Sleep(10 * time.Millisecond)
}

// readStream reads data from stream
func readStream(stream chan []Response) {

	// init counters
	mCtr := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: reader,
		Name:      "message",
	})

	bCtr := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: reader,
		Name:      "batch",
	})

	nCtr := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: reader,
		Name:      "no_messages_available",
	})

	// read from stream
	for {
		select {
		case responses := <-stream:

			// iterate over batch of responses
			for _, r := range responses {

				// read each message in a response
				for _, m := range r.Messages {

					// every message must be deleted
					deleteMessage(m)

					mCtr.Inc()
				}
			}

			bCtr.Inc()
		default:
			nCtr.Inc()
			time.Sleep(2 * time.Second)
		}
	}
}

// writeStream sends data to stream
func writeStream(stream chan []Response) {

	// init counter metric to count batch sent
	bCtr := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: writer,
		Name:      "batch",
	})

	fCtr := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: writer,
		Name:      "stream_full",
	})

	// loop forever
	for {

		// generate batch of responses
		var responses []Response
		for i := 0; i < responseBatchSize; i++ {

			// generate batch of messages
			var messages []Message
			for i := 0; i < msgBatchSize; i++ {

				// generate random payload
				payload := make([]byte, 1000)
				_, err := rand.Read(payload)
				if err != nil {
					fmt.Println("rand.Read error:", err)
					continue
				}

				m := Message{
					Data: payload,
				}

				messages = append(messages, m)
			}

			// create response
			response := Response{
				Messages: messages,
			}

			// add response to batch
			responses = append(responses, response)
		}

		// send batch to stream
	sendLoop:
		for {
			select {
			case stream <- responses:
				bCtr.Inc()
				break sendLoop
			default:
				fCtr.Inc()
				time.Sleep(3 * time.Second)
			}
		}
	}
}
