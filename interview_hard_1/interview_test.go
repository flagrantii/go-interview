package interview_hard_1

import (
	"context"
	"sync"
	"testing"
	"time"
)

// Test advanced concurrency patterns
func TestAdvancedConcurrency(t *testing.T) {
	t.Run("test rate limiter", func(t *testing.T) {
		requests := 10
		rateLimit := 2 // 2 requests per second
		start := time.Now()

		results := ProcessWithRateLimit(requests, rateLimit)

		duration := time.Since(start)
		expectedMinDuration := time.Duration(requests/rateLimit) * time.Second

		if duration < expectedMinDuration {
			t.Errorf("Rate limiting failed: took %v, expected at least %v", duration, expectedMinDuration)
		}

		if len(results) != requests {
			t.Errorf("Expected %d results, got %d", requests, len(results))
		}
	})

	t.Run("test context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		result, err := ProcessWithContext(ctx, 1000) // Should take longer than timeout
		if err == nil {
			t.Error("Expected context timeout error, got nil")
		}
		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})
}

// Test advanced data structures
func TestAdvancedDataStructures(t *testing.T) {
	t.Run("test concurrent safe cache", func(t *testing.T) {
		cache := NewThreadSafeCache()
		var wg sync.WaitGroup

		// Concurrent writes and reads
		for i := 0; i < 100; i++ {
			wg.Add(2)
			go func(val int) {
				defer wg.Done()
				cache.Set(val, val*2)
			}(i)
			go func(val int) {
				defer wg.Done()
				_, _ = cache.Get(val)
			}(i)
		}

		wg.Wait()
		if size := cache.Size(); size != 100 {
			t.Errorf("Expected cache size 100, got %d", size)
		}
	})

	t.Run("test priority queue", func(t *testing.T) {
		pq := NewPriorityQueue()
		items := []int{3, 1, 4, 1, 5, 9, 2, 6}
		for _, item := range items {
			pq.Push(item)
		}

		expected := []int{1, 1, 2, 3, 4, 5, 6, 9}
		for i, want := range expected {
			got := pq.Pop()
			if got != want {
				t.Errorf("item %d: got %d, want %d", i, got, want)
			}
		}
	})
}

// Test advanced error handling and recovery
func TestAdvancedErrorHandling(t *testing.T) {
	t.Run("test panic recovery", func(t *testing.T) {
		result, err := ExecuteWithRecover(func() interface{} {
			panic("expected panic")
		})

		if err == nil {
			t.Error("Expected error from panic, got nil")
		}
		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})

	t.Run("test error wrapping", func(t *testing.T) {
		_, err := ProcessWithErrorWrapping(-1)
		if err == nil {
			t.Error("Expected error for negative input")
		}

		var validationErr *ValidationError
		if !IsValidationError(err, &validationErr) {
			t.Error("Expected ValidationError type")
		}
	})
}

// Test advanced interfaces and type assertions
func TestAdvancedInterfaces(t *testing.T) {
	t.Run("test interface composition", func(t *testing.T) {
		var processor DataProcessor = &AdvancedProcessor{}

		if _, ok := processor.(DataValidator); !ok {
			t.Error("Expected processor to implement DataValidator")
		}

		if _, ok := processor.(DataTransformer); !ok {
			t.Error("Expected processor to implement DataTransformer")
		}
	})
}

// Test reflection and meta-programming
func TestReflection(t *testing.T) {
	t.Run("test struct field mapping", func(t *testing.T) {
		type TestStruct struct {
			Name  string `json:"name" validate:"required"`
			Age   int    `json:"age" validate:"min=0"`
			Email string `json:"email" validate:"email"`
		}

		ts := TestStruct{Name: "Test", Age: 25, Email: "test@example.com"}
		fields := GetStructFields(ts)

		expectedFields := map[string]string{
			"Name":  "name",
			"Age":   "age",
			"Email": "email",
		}

		for field, jsonTag := range expectedFields {
			if fields[field] != jsonTag {
				t.Errorf("Field %s: expected tag %s, got %s", field, jsonTag, fields[field])
			}
		}
	})
}

// Test advanced channel patterns
func TestAdvancedChannelPatterns(t *testing.T) {
	t.Run("test fan-out fan-in", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		workers := 3
		expected := 385 // sum of squares of numbers 1-10

		result := FanOutFanIn(input, workers)
		if result != expected {
			t.Errorf("Expected %d, got %d", expected, result)
		}
	})

	t.Run("test pipeline pattern", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5}
		expected := []int{4, 8, 12, 16, 20}

		result := Pipeline(numbers)
		for i, v := range expected {
			if result[i] != v {
				t.Errorf("At index %d: expected %d, got %d", i, v, result[i])
			}
		}
	})
}
