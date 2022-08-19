package worklist

type Entry struct {
	path string
}

type WorkList struct {
	jobs chan Entry
}

func (w *WorkList) AddJob(work Entry) {
	// Placing work data into the jobs channel
	w.jobs <- work
}

func (w *WorkList) Next() Entry {
	jobs := <-w.jobs
	return jobs
}

func New(bufferSize int) WorkList {
	return WorkList{make(chan Entry, bufferSize)}
}

// func to create a new job with the the Entry struct
func NewJob(path string) Entry {
	return Entry{path: path}
}

func (w *WorkList) Finalize(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		w.AddJob(Entry{""})
	}
}
