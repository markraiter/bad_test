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
	fileLimit            = 100 * 1024 * 1024 //TODO: FIX
	file                 = "file"
)

type TaskService struct {
	log *slog.Logger
}

func (ts *TaskService) FindValues(form *multipart.Form) (string, error) {
	const operation = "service.Task.Foo"

	log := ts.log.With(slog.String("operation", operation))

	log.Info("attempting to process file")

	if err := validateFile(form); err != nil {
		return "", fmt.Errorf("%s: %w", operation, err)
	}

	sourceFile := form.File[file][0]

	rFile, err := sourceFile.Open()
	if err != nil {
		return "", fmt.Errorf("%s: %w", operation, err)
	}
	defer rFile.Close()

	scanner := bufio.NewScanner(rFile)
	var numbers []int
	var wg sync.WaitGroup
	var max, min, sum int
	var median float64

	for scanner.Scan() {
		line := scanner.Text()
		nums := strings.Fields(line)

		for _, numStr := range nums {
			num, err := strconv.Atoi(numStr)
			if err != nil {
				fmt.Errorf("%s: %w", operation, err)
				continue
			}

			numbers = append(numbers, num)
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

	log.Info("file processed\nMax: %d\nMin: %d\nMedian: %.2f\nAverage: %.2f", max, min, median, float64(sum)/float64(len(numbers)))

	result := fmt.Sprintf("Max: %d\nMin: %d\nMedian: %.2f\nAverage: %.2f", max, min, median, float64(sum)/float64(len(numbers)))
	return result, nil
}

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

func findSum(numbers []int) int {
	sum := 0
	for _, num := range numbers {
		sum += num
	}
	return sum
}

// validateFile validates provided file.
//
// Supported extensions: .txt
// Supported number of items: 1
// Supported file size: 25MB
// If file is not valid, returns error.
func validateFile(form *multipart.Form) error {
	allowedFileExtentions := []string{"txt", "TXT"}

	if len(form.File[file]) < limitNumberItemsFile {
		return fmt.Errorf("error: %w", model.ErrNoFile)
	}

	if len(form.File[file]) > limitNumberItemsFile {
		return fmt.Errorf("error: %w", model.ErrToManyFiles)
	}

	sourceFile := form.File[file][0]

	if sourceFile == nil || sourceFile.Size > fileLimit || !isAllowedFileExtention(allowedFileExtentions, sourceFile.Filename) {
		return fmt.Errorf("error: %w", model.ErrFileExtensionViolation)
	}

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
