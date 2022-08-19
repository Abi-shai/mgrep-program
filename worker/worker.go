package worker

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Result struct {
	line       string
	lineNumber int
	path       string
}

type Results struct {
	inner []Result
}

func NewResult(line string, lineNumb int, path string) Result {
	return Result{line, lineNumb, path}
}

// Searching an finding a file based on a string value
func FindInFile(path string, find string) *Results {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Sorry, couldn't find the file !")
	}

	results := Results{make([]Result, 0)}

	scanner := bufio.NewScanner(file)
	lineNumber := 1

	for scanner.Scan() {

		// If the text file contains the string to search
		if strings.Contains(scanner.Text(), find) {

			// Set the a new result
			response := NewResult(scanner.Text(), lineNumber, path)

			// Append and update the results from the *Results struct
			results.inner = append(results.inner, response)
		}
		lineNumber += 1
	}
	if len(results.inner) == 0 {
		return nil
	} else {
		return &results
	}
}
