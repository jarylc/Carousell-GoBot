package messaging

import (
	"log"
	"math"
	"os"
	"sync"
	"time"
)

var debug = os.Getenv("DEBUG") == "1"

type queueItem struct {
	messager Messager
	message  string
}

var queue = &[]queueItem{}
var queueMutex sync.Mutex
var queueProcessing = false
var retry = &[]queueItem{}
var retryMutex sync.Mutex
var retryProcessing = false
var retryInterval = new(int8)

func popList(item *[]queueItem) {
	queueMutex.Lock()
	retryMutex.Lock()
	defer queueMutex.Unlock()
	defer retryMutex.Unlock()

	*item = append((*item)[:0], (*item)[1:]...)
}
func processList(item *[]queueItem) error {
	if debug {
		log.Println("Processing a messaging queue")
	}
	err := (*item)[0].messager.ProcessMessage((*item)[0].message)
	if err != nil {
		log.Println(err)
		addRetry((*item)[0]) // move to back of retry queue
		popList(item)
		startRetry()
		return err
	}
	popList(item)
	return nil
}

func startRetry() {
	if retryProcessing {
		return
	}
	retryProcessing = true

	retryMutex.Lock()
	defer retryMutex.Unlock()

	*retryInterval = 5

	if debug {
		log.Println("Processing all retry list")
	}
	go func(retry *[]queueItem) {
		for {
			<-time.After(time.Duration(*retryInterval) * time.Second) // every 5 seconds (increases by 5 every failure, reset on success)
			if len(*retry) > 0 {
				err := processList(retry)
				if err != nil {
					*retryInterval = int8(math.Min(float64(*retryInterval+5), 60))
					if debug {
						log.Printf("Retry interval increased to %ds\n", *retryInterval)
					}
				} else {
					*retryInterval = 5
					if debug {
						log.Println("Retry interval reset to 5s")
					}
				}
			} else {
				retryProcessing = false
				if debug {
					log.Println("Done processing retry list")
				}
				return
			}
		}
	}(retry)
}
func addRetry(item queueItem) {
	retryMutex.Lock()
	defer retryMutex.Unlock()

	if debug {
		log.Println("Adding item to messaging retry list")
	}
	*retry = append(*retry, item)
}

func startQueue() {
	if queueProcessing {
		return
	}
	queueProcessing = true

	if debug {
		log.Println("Processing all messaging queue")
	}
	go func(queue *[]queueItem) {
		for {
			<-time.After(time.Second) // every 1 second
			if len(*queue) > 0 {
				_ = processList(queue)
			} else {
				queueProcessing = false
				if debug {
					log.Println("Done processing messaging queue")
				}
				return
			}
		}
	}(queue)
}
func addQueue(item queueItem) {
	queueMutex.Lock()
	defer queueMutex.Unlock()

	if debug {
		log.Println("Adding item to messaging queue")
	}
	*queue = append(*queue, item)

	startQueue()
}
