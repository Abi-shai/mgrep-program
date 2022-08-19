package main

import (
	"fmt"
	"mgrep/worker"
	"mgrep/worklist"
	"os"
	"path/filepath"
	"sync"

	"github.com/alexflint/go-arg"
)

func discoverDirectories(workList *worklist.WorkList, path string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Read directory erro:", err)
		return
	}

	for _, entry := range entries {
		// If we encounter another directory, recurs into it
		if entry.IsDir() {
			nextPath := filepath.Join(path, entry.Name())

			// Recalling the func for the current directory
			discoverDirectories(workList, nextPath)
		} else {
			// When entry != directory
			// add the new file job to the worklist
			workList.AddJob(worklist.NewJob(filepath.Join(path, entry.Name())))
		}
	}
}

var args struct {
	// Using the arg config from the go-arg package
	// to make var have additionnal requirements
	searchTerm string `arg:"positionnal,required"`
	searchDir  string `arg:"positionnal"`
}

func main() {
	arg.MustParse(&args)

	var workersWaitGroup sync.WaitGroup

	worklist := worklist.New(100)

	resultsWorker := make(chan worker.Result, 100)

	numberOfWorkers := 10

	workersWaitGroup.Add(1)

	// Creating a go routine that handles
	// discovering directories
	go func() {
		defer workersWaitGroup.Done()

		discoverDirectories(&worklist, args.searchDir)
		worklist.Finalize(numberOfWorkers)
	}()

	for i := 0; i < numberOfWorkers; i++ {
		workersWaitGroup.Add(1)
		go func() {
			defer workersWaitGroup.Done()
			for {
				workEntry := worklist.Next()
				if workEntry.path != "" {
					workerResult := worker.FindInFile(workEntry.Path, args.searchTerm)

					if workerResult != nil {
						for _, result := range workerResult.inner {
							results <- result
						}
					}
				} else {
					return
				}
			}
		}()
	}

	blockWorkersWaitGroup := make(chan struct{})
	go func() {
		workersWaitGroup.Wait()
		close(blockWorkersWaitGroup)
	}()

	var displayWaitGroup sync.WaitGroup

	displayWaitGroup.Add(1)
	go func() {
		for {
			select {
			case results := <-results:
				fmt.Printf("%v[%v]:%v\n", r.path, r.lineNumber, r.line)
			case <-blockWorkersWaitGroup:
				if len(results) == 0 {
					displayWaitGroup.Done()
					return
				}
			}
		}
	}()
	displayWaitGroup.Wait()
}
