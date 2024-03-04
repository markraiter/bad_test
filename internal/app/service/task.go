package service

import (
	"bufio"
	"fmt"
	"log/slog"
	"mime/multipart"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/markraiter/bad_test/internal/model"
)

const (
	nameLimit            = 2
	limitNumberItemsFile = 1
	file                 = "file"
)

type TaskService struct {
	log *slog.Logger
}

// FindValues finds min, max, median, average, max increasing sequential and max decreasing sequential of the numbers in the file.
//
// If file is not valid, returns error.
// If file is empty, returns error.
// If file is valid, returns min, max, median, average, max increasing sequential and max decreasing sequential of the numbers.
func (ts *TaskService) FindValues(form *multipart.Form) (*model.TaskResult, error) {
	const operation = "service.Task.FindValues"

	log := ts.log.With(slog.String("operation", operation))

	log.Info("attempting to process file")

	if err := ts.validateFile(form); err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	sourceFile := form.File[file][0]

	rFile, err := sourceFile.Open()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	defer rFile.Close()

	scanner := bufio.NewScanner(rFile)
	var numbers []int
	var wg sync.WaitGroup
	var max, min, sum int
	var median float64
	var prevNum int
	var increasingSeq, decreasingSeq []int
	var result model.TaskResult

	for scanner.Scan() {
		line := scanner.Text()
		nums := strings.Fields(line)

		for _, numStr := range nums {
			num, err := strconv.Atoi(numStr)
			if err != nil {
				log.Debug("error converting string to int", model.Err(err))
				continue
			}

			numbers = append(numbers, num)

			if num > prevNum {
				increasingSeq = append(increasingSeq, num)
			} else {
				if len(increasingSeq) > len(result.MaxIncreasingSeq) {
					result.MaxIncreasingSeq = increasingSeq
				}
				increasingSeq = []int{num}
			}

			if num < prevNum {
				decreasingSeq = append(decreasingSeq, num)
			} else {
				if len(decreasingSeq) > len(result.MaxDecreasingSeq) {
					result.MaxDecreasingSeq = decreasingSeq
				}
				decreasingSeq = []int{num}
			}

			prevNum = num
		}
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		findMinMax(numbers, &max, &min)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		median = findMedian(numbers)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		sum = findSum(numbers)
	}()

	wg.Wait()

	result.Max = max
	result.Min = min
	result.Median = median
	result.Avg = float64(sum) / float64(len(numbers))

	log.Info("file processed")

	return &result, nil
}

// validateFile validates provided file.
//
// Supported extensions: .txt
// Supported number of items: 1
// Supported file size: 100MB
// If file is not valid, returns error.
func (ts *TaskService) validateFile(form *multipart.Form) error {
	const operation = "service.Task.validateFile"

	log := ts.log.With(slog.String("operation", operation))

	log.Info("validating file")

	allowedFileExtentions := []string{"txt", "TXT"}

	if len(form.File[file]) < limitNumberItemsFile {
		log.Error("no file", model.Err(model.ErrNoFile))

		return fmt.Errorf("error: %w", model.ErrNoFile)
	}

	if len(form.File[file]) > limitNumberItemsFile {
		log.Error("too many files", model.Err(model.ErrToManyFiles))

		return fmt.Errorf("error: %w", model.ErrToManyFiles)
	}

	sourceFile := form.File[file][0]

	if sourceFile == nil || !isAllowedFileExtention(allowedFileExtentions, sourceFile.Filename) {
		log.Error("file extension violation", model.Err(model.ErrFileExtensionViolation))

		return fmt.Errorf("error: %w", model.ErrFileExtensionViolation)
	}

	log.Info("file is valid")

	return nil
}

// isAllowedFileExtention checks if file extension is allowed.
//
// Allowed extensions: .txt
func isAllowedFileExtention(allowedList []string, fileName string) bool {
	nameParts := strings.Split(fileName, ".")

	fileExt := nameParts[len(nameParts)-1]
	for _, i := range allowedList {
		if i == fileExt {
			return true
		}
	}

	return false
}

// findMedian finds median of the numbers.
//
// If numbers length is even, returns average of two middle numbers.
func findMedian(numbers []int) float64 {
	sorted := make([]int, len(numbers))

	copy(sorted, numbers)

	sort.Ints(sorted)

	n := len(sorted)

	if n%2 == 0 {
		return float64(sorted[n/2-1]+sorted[n/2]) / 2.0
	}

	return float64(sorted[n/2])
}

// findSum finds sum of the numbers.
//
// If numbers length is even, returns average of two middle numbers.
func findSum(numbers []int) int {
	sum := 0

	for _, num := range numbers {
		sum += num
	}

	return sum
}

// findMinMax finds min and max of the numbers.
//
// If numbers length is even, returns average of two middle numbers.
func findMinMax(numbers []int, max, min *int) {
	*max = numbers[0]
	*min = numbers[0]

	for _, num := range numbers {
		if num > *max {
			*max = num
		}

		if num < *min {
			*min = num
		}
	}
}
