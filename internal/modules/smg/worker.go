package smg

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// Worker analyzes the links in the site content and returns the result to the coordinator.
type Worker struct {
	name                  string
	busy                  bool
	mutex                 *sync.RWMutex
	listenChan            chan *workerPayload
	coordinatorSenderChan chan<- *coordinatorPayloads
}

// NewWorker returns a new worker under given coordinator
func NewWorker(name string, coordinatorSender chan<- *coordinatorPayloads) *Worker {
	return &Worker{
		name:                  name,
		mutex:                 &sync.RWMutex{},
		listenChan:            make(chan *workerPayload),
		coordinatorSenderChan: coordinatorSender,
	}
}

// Busy returns the worker's status
func (w *Worker) Busy() bool {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.busy
}

// digSite extract links from the analyzed page
func digSite(depth int, u *url.URL) (wd *coordinatorPayload) {
	resp, err := http.DefaultClient.Get(u.String())
	if err != nil {
		return &coordinatorPayload{
			depth:     depth + 1,
			sourceURL: u,
			err:       err,
		}
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &coordinatorPayload{
			depth:     depth + 1,
			sourceURL: u,
			err:       fmt.Errorf("url responsed with code %d", resp.StatusCode),
		}
	}

	ct := resp.Header.Get("Content-type")
	if ct != "" && !strings.Contains(ct, "text/html") {
		return &coordinatorPayload{
			depth:     depth + 1,
			sourceURL: u,
			err:       fmt.Errorf("unknown content type: %s", ct),
		}
	}

	s, iu := ParseHTML(u, resp.Body)
	return &coordinatorPayload{
		depth:       depth + 1,
		sourceURL:   u,
		urls:        s,
		invalidURLs: iu,
	}
}

// digSites extract links from the analyzed pages
func digSites(depth int, urls []*url.URL) (wds []*coordinatorPayload) {
	for _, u := range urls {
		wds = append(wds, digSite(depth, u))
	}
	return wds
}

// Run initialize the worker
func (w *Worker) Run(ctx context.Context) {
	log.Printf("Starting worker %s...\n", w.name)
	for {
		select {
		case <-ctx.Done():
			return
		case message := <-w.listenChan:
			w.busy = true
			log.Printf("Digging urls(%d) from depth %d\n", len(message.urls), message.currentDepth)
			wds := digSites(message.currentDepth, message.urls)
			got := make(chan bool)
			w.coordinatorSenderChan <- &coordinatorPayloads{
				worker: w.name,
				got:    got,
				dumps:  wds,
			}
			<-got
			w.busy = false
		}
	}
}
