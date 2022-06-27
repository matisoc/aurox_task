package smg

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"testing"
)

func Test_availableWorkers(t *testing.T) {
	coordinator := &Coordinator{
		workers: []*Worker{
			{
				name:  "w1",
				busy:  true,
				mutex: &sync.RWMutex{},
			},

			{
				name:  "w2",
				mutex: &sync.RWMutex{},
			},

			{
				name:  "w3",
				mutex: &sync.RWMutex{},
			},
			{
				name:  "w4",
				busy:  true,
				mutex: &sync.RWMutex{},
			},
		},
	}

	expectedAvailableWorkers := map[string]bool{
		"w2": true,
		"w3": true,
	}

	availableWorkers := getAvailableWorkers(coordinator)
	if len(availableWorkers) != len(expectedAvailableWorkers) {
		t.Fatalf("expected available workers %d but got %d", len(expectedAvailableWorkers), len(availableWorkers))
	}

	for _, w := range availableWorkers {
		_, ok := expectedAvailableWorkers[w.name]
		if !ok {
			t.Fatalf("expected %s to be available but busy\n", w.name)
		}
	}
}

func Test_distributePayload(t *testing.T) {
	data := []struct {
		workers int
		busy    int
		urls    int
	}{
		{
			workers: 4,
			busy:    4,
		},

		{
			workers: 5,
			urls:    4,
		},

		{
			workers: 5,
			busy:    1,
			urls:    4,
		},

		{
			workers: 5,
			urls:    10,
		},
	}

	workersBuild := func(c *Coordinator, total, busy int) (workers []*Worker) {
		for i := 0; i < total; i++ {
			workers = append(workers, NewWorker(fmt.Sprintf("worker %d", i), c.submitCh))
		}

		for i := 0; i < busy; i++ {
			workers[i].busy = true
		}

		return workers
	}

	urlBuild := func(total int) (urls []*url.URL) {
		for i := 0; i < total; i++ {
			u, _ := url.Parse(fmt.Sprintf("https://google.com/%d", i))
			urls = append(urls, u)
		}

		return urls
	}

	for _, test := range data {
		baseURL, _ := url.Parse("https://google.com")
		coordinator := NewCoordinator(context.Background(), baseURL, 1, 1)
		coordinator.workers = workersBuild(coordinator, test.workers, test.busy)
		urls := urlBuild(test.urls)

		testCh := make(chan int)
		for _, w := range coordinator.workers {
			go func(w *Worker) {
				for wp := range w.listenChan {
					testCh <- len(wp.urls)
				}
			}(w)
		}
		distributePayload(coordinator, 1, urls)

		count := 0
		busyWorkers := test.urls
		if busyWorkers > (test.workers - test.busy) {
			busyWorkers = test.workers - test.busy
		}
		for i := 0; i < busyWorkers; i++ {
			count += <-testCh
		}

		if test.urls != count {
			t.Fatalf("expected %d urls but got %d", test.urls, count)
		}

	}
}
