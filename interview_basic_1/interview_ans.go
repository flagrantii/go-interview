package interview_basic_1

import (
	"errors"
	"strings"
	"sync"
)

func SumNumbers(numbers []int) int {
	sum := 0
	for _, num := range numbers {
		sum += num
	}

	return sum
}

func CountFrequency(objs []string) map[string]int {
	ans := make(map[string]int)
	for _, obj := range objs {
		ans[obj]++
	}
	return ans
}

func ReverseSlice(slice []int) []int {
	temp := []int{}
	for i := len(slice) - 1; i >= 0; i-- {
		temp = append(temp, slice[i])
	}
	return temp
}

type Rectangle struct {
	Width  int
	Height int
}

func (r *Rectangle) Area() float64 {
	return float64(r.Height * r.Width)
}

func (r *Rectangle) Perimeter() float64 {
	return float64(2 * (r.Height + r.Width))
}

func SwapValues(x *int, y *int) {
	*x, *y = *y, *x
}

func ParallelSquare(numbers []int) []int {
	result := make([]int, len(numbers))
	var wg sync.WaitGroup

	for i, num := range numbers {
		wg.Add(1)
		go func(idx, value int) {
			defer wg.Done()
			result[idx] = value * value
		}(i, num)
	}

	wg.Wait()
	return result
}

func ProcessWithWorkerPool(tasks []int, workers int) int {
	jobs := make(chan int, len(tasks))
	result := make(chan int, len(tasks))
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(jobs, result, &wg)
	}

	for _, task := range tasks {
		jobs <- task
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(result)
	}()

	ans := 0
	for reslt := range result {
		ans += reslt
	}
	return ans
}

func worker(jobs <-chan int, result chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for num := range jobs {
		result <- num
	}
}

func DivideNumbers(a int, b int) (float64, error) {
	if b == 0 {
		return 0, errors.New("b cannot be 0")
	}
	return float64(a) / float64(b), nil
}

type CustomString string

func (cs CustomString) ToUpper() CustomString {
	return CustomString(strings.ToUpper(string(cs)))
}
