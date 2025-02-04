package interview_basic_2

import (
	"context"
	"sync"
	"testing"
	"time"
)

// Test concurrent binary tree operations
func TestConcurrentBinaryTree(t *testing.T) {
	t.Run("test concurrent insertions and searches", func(t *testing.T) {
		tree := NewThreadSafeBST()
		numbers := []int{5, 3, 7, 1, 9, 4, 6, 8, 2, 0}
		var wg sync.WaitGroup

		// Concurrent insertions
		for _, num := range numbers {
			wg.Add(1)
			go func(val int) {
				defer wg.Done()
				tree.Insert(val)
			}(num)
		}
		wg.Wait()

		// Verify all numbers are present
		for _, num := range numbers {
			if !tree.Search(num) {
				t.Errorf("Number %d not found in tree", num)
			}
		}
	})
}

// Test custom generic cache with expiration
func TestGenericCache(t *testing.T) {
	t.Run("test cache operations with expiration", func(t *testing.T) {
		cache := NewCache[string, int](100 * time.Millisecond)

		cache.Set("one", 1)
		cache.Set("two", 2)

		// Test immediate retrieval
		if val, ok := cache.Get("one"); !ok || val != 1 {
			t.Errorf("Expected 1, got %v", val)
		}

		// Wait for expiration
		time.Sleep(150 * time.Millisecond)

		// Should be expired
		if _, ok := cache.Get("one"); ok {
			t.Error("Cache entry should have expired")
		}
	})
}

// Test custom rate limiter with sliding window
func TestSlidingWindowRateLimiter(t *testing.T) {
	t.Run("test rate limiting", func(t *testing.T) {
		limiter := NewRateLimiter(3, time.Second) // 3 requests per second

		// First three requests should succeed
		for i := 0; i < 3; i++ {
			if !limiter.Allow() {
				t.Errorf("Request %d should be allowed", i)
			}
		}

		// Fourth request should be denied
		if limiter.Allow() {
			t.Error("Fourth request should be denied")
		}

		// Wait and try again
		time.Sleep(time.Second)
		if !limiter.Allow() {
			t.Error("Request should be allowed after window reset")
		}
	})
}

// Test custom event bus with middleware
func TestEventBus(t *testing.T) {
	t.Run("test event publishing and subscription", func(t *testing.T) {
		bus := NewEventBus()
		received := make(chan string, 1)

		// Add middleware
		bus.Use(func(evt Event) Event {
			evt.Payload = evt.Payload.(string) + "_modified"
			return evt
		})

		// Subscribe
		bus.Subscribe("test_event", func(evt Event) {
			received <- evt.Payload.(string)
		})

		// Publish
		bus.Publish("test_event", "hello")

		select {
		case msg := <-received:
			if msg != "hello_modified" {
				t.Errorf("Expected 'hello_modified', got %s", msg)
			}
		case <-time.After(time.Second):
			t.Error("Timeout waiting for event")
		}
	})
}

// Test custom scheduler with priority
func TestPriorityScheduler(t *testing.T) {
	t.Run("test task scheduling with priorities", func(t *testing.T) {
		scheduler := NewScheduler(context.Background())
		results := make(chan int, 3)

		scheduler.Schedule(Task{
			Priority: 3,
			Execute: func() {
				results <- 3
			},
		})

		scheduler.Schedule(Task{
			Priority: 1,
			Execute: func() {
				results <- 1
			},
		})

		scheduler.Schedule(Task{
			Priority: 2,
			Execute: func() {
				results <- 2
			},
		})

		expected := []int{1, 2, 3}
		for i, want := range expected {
			select {
			case got := <-results:
				if got != want {
					t.Errorf("Task %d: expected priority %d, got %d", i, want, got)
				}
			case <-time.After(time.Second):
				t.Error("Timeout waiting for task execution")
			}
		}
	})
}

// Test custom retry mechanism with backoff
func TestRetryMechanism(t *testing.T) {
	t.Run("test retry with exponential backoff", func(t *testing.T) {
		attempts := 0
		start := time.Now()

		err := RetryWithBackoff(func() error {
			attempts++
			if attempts < 3 {
				return ErrTemporary
			}
			return nil
		}, 3, 50*time.Millisecond)

		duration := time.Since(start)
		if err != nil {
			t.Errorf("Expected success, got error: %v", err)
		}

		if attempts != 3 {
			t.Errorf("Expected 3 attempts, got %d", attempts)
		}

		// Check if backoff timing is reasonable
		if duration < 150*time.Millisecond {
			t.Error("Backoff duration too short")
		}
	})
}

// Test string manipulation
func TestStringManipulation(t *testing.T) {
	t.Run("test palindrome check", func(t *testing.T) {
		tests := []struct {
			input    string
			expected bool
		}{
			{"radar", true},
			{"hello", false},
			{"A man a plan a canal Panama", true},
			{"race a car", false},
			{"", true},
		}

		for _, tt := range tests {
			result := IsPalindrome(tt.input)
			if result != tt.expected {
				t.Errorf("IsPalindrome(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		}
	})
}

// Test slice operations
func TestSliceOperations(t *testing.T) {
	t.Run("test remove duplicates", func(t *testing.T) {
		input := []int{1, 1, 2, 2, 3, 4, 4, 5}
		expected := []int{1, 2, 3, 4, 5}

		result := RemoveDuplicates(input)
		if len(result) != len(expected) {
			t.Errorf("RemoveDuplicates() length = %v, want %v", len(result), len(expected))
		}
		for i, v := range expected {
			if result[i] != v {
				t.Errorf("RemoveDuplicates() at index %d = %v, want %v", i, result[i], v)
			}
		}
	})
}

// Interface for geometric shapes
type Shape2D interface {
	Area() float64
	Perimeter() float64
	Scale(factor float64)
}

// Test struct and interface implementation
func TestShapeOperations(t *testing.T) {
	t.Run("test circle operations", func(t *testing.T) {
		circle := &Circle{Radius: 5}
		expectedArea := 78.53981633974483
		expectedPerimeter := 31.41592653589793

		if area := circle.Area(); !almostEqual(area, expectedArea) {
			t.Errorf("Circle.Area() = %v, want %v", area, expectedArea)
		}

		if perimeter := circle.Perimeter(); !almostEqual(perimeter, expectedPerimeter) {
			t.Errorf("Circle.Perimeter() = %v, want %v", perimeter, expectedPerimeter)
		}

		circle.Scale(2)
		if radius := circle.Radius; radius != 10 {
			t.Errorf("Circle.Scale(2) radius = %v, want 10", radius)
		}
	})
}

// Helper function for float comparison
func almostEqual(a, b float64) bool {
	const epsilon = 1e-10
	return abs(a-b) <= epsilon
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// Test map operations
func TestMapOperations(t *testing.T) {
	t.Run("test word frequency", func(t *testing.T) {
		input := "the quick brown fox jumps over the lazy dog"
		expected := map[string]int{
			"the":   2,
			"quick": 1,
			"brown": 1,
			"fox":   1,
			"jumps": 1,
			"over":  1,
			"lazy":  1,
			"dog":   1,
		}

		result := WordFrequency(input)
		for k, v := range expected {
			if result[k] != v {
				t.Errorf("WordFrequency() for word %q = %v, want %v", k, result[k], v)
			}
		}
	})
}

// Test error handling
func TestErrorHandling(t *testing.T) {
	t.Run("test array index validation", func(t *testing.T) {
		arr := []int{1, 2, 3, 4, 5}
		tests := []struct {
			index    int
			wantErr  bool
			expected int
		}{
			{2, false, 3},
			{-1, true, 0},
			{5, true, 0},
			{0, false, 1},
		}

		for _, tt := range tests {
			result, err := GetElement(arr, tt.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetElement(arr, %d) error = %v, wantErr %v", tt.index, err, tt.wantErr)
				continue
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("GetElement(arr, %d) = %v, want %v", tt.index, result, tt.expected)
			}
		}
	})
}

// Add this test function
func TestConcurrentCounter(t *testing.T) {
	t.Run("test concurrent counter operations", func(t *testing.T) {
		counter := NewThreadSafeCounter()
		var wg sync.WaitGroup
		iterations := 1000
		goroutines := 10

		// Test concurrent increments
		for i := 0; i < goroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < iterations; j++ {
					counter.Increment()
				}
			}()
		}

		wg.Wait()
		expected := goroutines * iterations
		if value := counter.GetValue(); value != expected {
			t.Errorf("Counter value = %d, want %d", value, expected)
		}

		// Test concurrent get and set
		wg.Add(2)
		go func() {
			defer wg.Done()
			counter.SetValue(42)
		}()
		go func() {
			defer wg.Done()
			_ = counter.GetValue()
		}()
		wg.Wait()
	})
}
