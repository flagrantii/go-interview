package interview_hard_1

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"time"
)

// Rate Limiter Implementation
func ProcessWithRateLimit(requests, rateLimit int) []int {
	results := make([]int, requests)
	ticker := time.NewTicker(time.Second / time.Duration(rateLimit))
	defer ticker.Stop()

	for i := 0; i < requests; i++ {
		<-ticker.C
		results[i] = i
	}
	return results
}

// Context Implementation
func ProcessWithContext(ctx context.Context, duration int) ([]int, error) {
	resultChan := make(chan []int, 1)

	go func() {
		time.Sleep(time.Duration(duration) * time.Millisecond)
		resultChan <- []int{1, 2, 3}
	}()

	select {
	case result := <-resultChan:
		return result, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Thread-safe Cache Implementation
type ThreadSafeCache struct {
	sync.RWMutex
	data map[interface{}]interface{}
}

func NewThreadSafeCache() *ThreadSafeCache {
	return &ThreadSafeCache{
		data: make(map[interface{}]interface{}),
	}
}

func (c *ThreadSafeCache) Set(key, value interface{}) {
	c.Lock()
	defer c.Unlock()
	c.data[key] = value
}

func (c *ThreadSafeCache) Get(key interface{}) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()
	val, ok := c.data[key]
	return val, ok
}

func (c *ThreadSafeCache) Size() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.data)
}

// Priority Queue Implementation
type PriorityQueue struct {
	items []int
}

func NewPriorityQueue() *PriorityQueue {
	return &PriorityQueue{items: make([]int, 0)}
}

func (pq *PriorityQueue) Push(item int) {
	pq.items = append(pq.items, item)
	pq.bubbleUp(len(pq.items) - 1)
}

func (pq *PriorityQueue) Pop() int {
	if len(pq.items) == 0 {
		return 0
	}
	item := pq.items[0]
	lastIdx := len(pq.items) - 1
	pq.items[0] = pq.items[lastIdx]
	pq.items = pq.items[:lastIdx]
	if len(pq.items) > 0 {
		pq.bubbleDown(0)
	}
	return item
}

func (pq *PriorityQueue) bubbleUp(index int) {
	for index > 0 {
		parent := (index - 1) / 2
		if pq.items[parent] > pq.items[index] {
			pq.items[parent], pq.items[index] = pq.items[index], pq.items[parent]
			index = parent
		} else {
			break
		}
	}
}

func (pq *PriorityQueue) bubbleDown(index int) {
	for {
		smallest := index
		left := 2*index + 1
		right := 2*index + 2

		if left < len(pq.items) && pq.items[left] < pq.items[smallest] {
			smallest = left
		}
		if right < len(pq.items) && pq.items[right] < pq.items[smallest] {
			smallest = right
		}

		if smallest == index {
			break
		}
		pq.items[index], pq.items[smallest] = pq.items[smallest], pq.items[index]
		index = smallest
	}
}

// Error Handling Implementation
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func ExecuteWithRecover(fn func() interface{}) (result interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("recovered from panic: " + r.(string))
		}
	}()
	return fn(), nil
}

func ProcessWithErrorWrapping(val int) (int, error) {
	if val < 0 {
		return 0, &ValidationError{Message: "value must be non-negative"}
	}
	return val * 2, nil
}

func IsValidationError(err error, target **ValidationError) bool {
	return errors.As(err, target)
}

// Interface Implementation
type DataProcessor interface {
	Process(data []byte) ([]byte, error)
}

type DataValidator interface {
	Validate(data []byte) error
}

type DataTransformer interface {
	Transform(data []byte) []byte
}

type AdvancedProcessor struct{}

func (ap *AdvancedProcessor) Process(data []byte) ([]byte, error) {
	if err := ap.Validate(data); err != nil {
		return nil, err
	}
	return ap.Transform(data), nil
}

func (ap *AdvancedProcessor) Validate(data []byte) error {
	if len(data) == 0 {
		return errors.New("empty data")
	}
	return nil
}

func (ap *AdvancedProcessor) Transform(data []byte) []byte {
	result := make([]byte, len(data))
	copy(result, data)
	return result
}

// Reflection Implementation
func GetStructFields(obj interface{}) map[string]string {
	result := make(map[string]string)
	t := reflect.TypeOf(obj)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if jsonTag := field.Tag.Get("json"); jsonTag != "" {
			result[field.Name] = jsonTag
		}
	}
	return result
}

// Channel Patterns Implementation
func FanOutFanIn(input []int, workers int) int {
	jobs := make(chan int, len(input))
	results := make(chan int, len(input))
	var wg sync.WaitGroup

	// Fan Out
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for num := range jobs {
				results <- num * num
			}
		}()
	}

	// Send jobs
	for _, num := range input {
		jobs <- num
	}
	close(jobs)

	// Wait and close results
	go func() {
		wg.Wait()
		close(results)
	}()

	// Fan In
	sum := 0
	for result := range results {
		sum += result
	}
	return sum
}

func Pipeline(numbers []int) []int {
	// Stage 1: Multiply by 2
	multiply := func(in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			for n := range in {
				out <- n * 2
			}
			close(out)
		}()
		return out
	}

	// Stage 2: Multiply by 2 again
	multiplyAgain := func(in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			for n := range in {
				out <- n * 2
			}
			close(out)
		}()
		return out
	}

	// Create input channel
	input := make(chan int)
	go func() {
		for _, n := range numbers {
			input <- n
		}
		close(input)
	}()

	// Connect the pipeline
	stage1 := multiply(input)
	stage2 := multiplyAgain(stage1)

	// Collect results
	result := make([]int, 0, len(numbers))
	for n := range stage2 {
		result = append(result, n)
	}
	return result
}
