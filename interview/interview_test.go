package interview

import (
	"testing"
)

// Interface definition for testing interface implementation
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Test basic syntax, operators, and control structures
func TestBasicOperations(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected int
	}{
		{
			name:     "sum of numbers",
			input:    []int{1, 2, 3, 4, 5},
			expected: 15,
		},
		{
			name:     "sum of negative numbers",
			input:    []int{-1, -2, -3},
			expected: -6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SumNumbers(tt.input)
			if result != tt.expected {
				t.Errorf("SumNumbers() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test slices and maps operations
func TestDataStructures(t *testing.T) {
	t.Run("test map operations", func(t *testing.T) {
		input := []string{"apple", "banana", "apple", "cherry", "banana"}
		expected := map[string]int{
			"apple":  2,
			"banana": 2,
			"cherry": 1,
		}

		result := CountFrequency(input)
		for k, v := range expected {
			if result[k] != v {
				t.Errorf("CountFrequency() for %s = %v, want %v", k, result[k], v)
			}
		}
	})

	t.Run("test slice operations", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		expected := []int{5, 4, 3, 2, 1}

		result := ReverseSlice(input)
		for i := range expected {
			if result[i] != expected[i] {
				t.Errorf("ReverseSlice() at index %d = %v, want %v", i, result[i], expected[i])
			}
		}
	})
}

// Test struct and method implementation
func TestStructsAndMethods(t *testing.T) {
	t.Run("test rectangle implementation", func(t *testing.T) {
		rect := &Rectangle{Width: 5, Height: 3}
		expectedArea := 15.0
		expectedPerimeter := 16.0

		if area := rect.Area(); area != expectedArea {
			t.Errorf("Rectangle.Area() = %v, want %v", area, expectedArea)
		}

		if perimeter := rect.Perimeter(); perimeter != expectedPerimeter {
			t.Errorf("Rectangle.Perimeter() = %v, want %v", perimeter, expectedPerimeter)
		}
	})
}

// Test pointer operations
func TestPointerOperations(t *testing.T) {
	t.Run("test swap function", func(t *testing.T) {
		x, y := 5, 10
		SwapValues(&x, &y)

		if x != 10 || y != 5 {
			t.Errorf("SwapValues() failed, got x=%v, y=%v, want x=10, y=5", x, y)
		}
	})
}

// Test goroutines and channels
func TestConcurrency(t *testing.T) {
	t.Run("test parallel processing", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5}
		expected := []int{1, 4, 9, 16, 25}

		result := ParallelSquare(numbers)
		for i := range expected {
			if result[i] != expected[i] {
				t.Errorf("ParallelSquare() at index %d = %v, want %v", i, result[i], expected[i])
			}
		}
	})

	t.Run("test worker pool", func(t *testing.T) {
		tasks := []int{1, 2, 3, 4, 5}
		expected := 15

		result := ProcessWithWorkerPool(tasks, 3)
		if result != expected {
			t.Errorf("ProcessWithWorkerPool() = %v, want %v", result, expected)
		}
	})
}

// Test error handling
func TestErrorHandling(t *testing.T) {
	t.Run("test division with error", func(t *testing.T) {
		tests := []struct {
			a, b     int
			expected float64
			wantErr  bool
		}{
			{10, 2, 5.0, false},
			{10, 0, 0, true},
		}

		for _, tt := range tests {
			result, err := DivideNumbers(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("DivideNumbers() error = %v, wantErr %v", err, tt.wantErr)
				continue
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("DivideNumbers() = %v, want %v", result, tt.expected)
			}
		}
	})
}

// Test custom type and interface implementation
func TestCustomTypes(t *testing.T) {
	t.Run("test custom string type", func(t *testing.T) {
		var cs CustomString = "hello"
		expected := "HELLO"

		if result := cs.ToUpper(); result != CustomString(expected) {
			t.Errorf("CustomString.ToUpper() = %v, want %v", result, expected)
		}
	})
}
