package interview_basic_2

import (
	"context"
	"errors"
	"math"
	"strings"
	"sync"
	"time"
)

// Thread-safe Binary Search Tree
type Node struct {
	Value int
	Left  *Node
	Right *Node
}

type ThreadSafeBST struct {
	root  *Node
	mutex sync.RWMutex
}

func NewThreadSafeBST() *ThreadSafeBST {
	return &ThreadSafeBST{}
}

func (t *ThreadSafeBST) Insert(value int) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if t.root == nil {
		t.root = &Node{Value: value}
		return
	}

	current := t.root
	for {
		if value < current.Value {
			if current.Left == nil {
				current.Left = &Node{Value: value}
				return
			}
			current = current.Left
		} else {
			if current.Right == nil {
				current.Right = &Node{Value: value}
				return
			}
			current = current.Right
		}
	}
}

func (t *ThreadSafeBST) Search(value int) bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	current := t.root
	for current != nil {
		if current.Value == value {
			return true
		}
		if value < current.Value {
			current = current.Left
		} else {
			current = current.Right
		}
	}
	return false
}

// Generic Cache with Expiration
type CacheEntry[V any] struct {
	value      V
	expiration time.Time
}

type Cache[K comparable, V any] struct {
	items map[K]CacheEntry[V]
	mu    sync.RWMutex
	ttl   time.Duration
}

func NewCache[K comparable, V any](ttl time.Duration) *Cache[K, V] {
	cache := &Cache[K, V]{
		items: make(map[K]CacheEntry[V]),
		ttl:   ttl,
	}
	go cache.cleanup()
	return cache
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = CacheEntry[V]{
		value:      value,
		expiration: time.Now().Add(c.ttl),
	}
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.items[key]
	if !exists || time.Now().After(entry.expiration) {
		var zero V
		return zero, false
	}
	return entry.value, true
}

func (c *Cache[K, V]) cleanup() {
	ticker := time.NewTicker(c.ttl)
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for k, v := range c.items {
			if now.After(v.expiration) {
				delete(c.items, k)
			}
		}
		c.mu.Unlock()
	}
}

// Sliding Window Rate Limiter
type RateLimiter struct {
	requests []time.Time
	limit    int
	window   time.Duration
	mu       sync.Mutex
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make([]time.Time, 0),
		limit:    limit,
		window:   window,
	}
}

func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	// Remove old requests
	valid := 0
	for _, t := range rl.requests {
		if t.After(windowStart) {
			rl.requests[valid] = t
			valid++
		}
	}
	rl.requests = rl.requests[:valid]

	if len(rl.requests) >= rl.limit {
		return false
	}

	rl.requests = append(rl.requests, now)
	return true
}

// Event Bus with Middleware
type Event struct {
	Type    string
	Payload interface{}
}

type EventHandler func(Event)
type Middleware func(Event) Event

type EventBus struct {
	subscribers map[string][]EventHandler
	middleware  []Middleware
	mu          sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]EventHandler),
		middleware:  make([]Middleware, 0),
	}
}

func (eb *EventBus) Subscribe(eventType string, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.subscribers[eventType] = append(eb.subscribers[eventType], handler)
}

func (eb *EventBus) Publish(eventType string, payload interface{}) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	event := Event{Type: eventType, Payload: payload}
	for _, middleware := range eb.middleware {
		event = middleware(event)
	}

	if handlers, exists := eb.subscribers[eventType]; exists {
		for _, handler := range handlers {
			go handler(event)
		}
	}
}

func (eb *EventBus) Use(middleware Middleware) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.middleware = append(eb.middleware, middleware)
}

// Priority Scheduler
type Task struct {
	Priority int
	Execute  func()
}

type Scheduler struct {
	tasks  []Task
	mu     sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc
}

func NewScheduler(ctx context.Context) *Scheduler {
	ctx, cancel := context.WithCancel(ctx)
	s := &Scheduler{
		tasks:  make([]Task, 0),
		ctx:    ctx,
		cancel: cancel,
	}
	go s.run()
	return s
}

func (s *Scheduler) Schedule(task Task) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Insert task in priority order
	pos := 0
	for i, t := range s.tasks {
		if task.Priority < t.Priority {
			pos = i
			break
		}
		pos = i + 1
	}

	s.tasks = append(s.tasks, Task{})
	copy(s.tasks[pos+1:], s.tasks[pos:])
	s.tasks[pos] = task
}

func (s *Scheduler) run() {
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			s.mu.Lock()
			if len(s.tasks) > 0 {
				task := s.tasks[0]
				s.tasks = s.tasks[1:]
				s.mu.Unlock()
				task.Execute()
			} else {
				s.mu.Unlock()
				time.Sleep(10 * time.Millisecond)
			}
		}
	}
}

// Retry Mechanism with Exponential Backoff
var ErrTemporary = errors.New("temporary error")

func RetryWithBackoff(fn func() error, maxAttempts int, initialDelay time.Duration) error {
	var err error
	delay := initialDelay

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err = fn()
		if err == nil {
			return nil
		}

		if attempt == maxAttempts {
			break
		}

		time.Sleep(delay)
		delay *= 2 // Exponential backoff
	}

	return err
}

// String Manipulation
func IsPalindrome(s string) bool {
	// Convert to lowercase and remove non-alphanumeric characters
	s = strings.ToLower(s)
	cleaned := ""
	for _, char := range s {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') {
			cleaned += string(char)
		}
	}

	// Check palindrome
	for i := 0; i < len(cleaned)/2; i++ {
		if cleaned[i] != cleaned[len(cleaned)-1-i] {
			return false
		}
	}
	return true
}

// Slice Operations
func RemoveDuplicates(nums []int) []int {
	if len(nums) == 0 {
		return nums
	}

	seen := make(map[int]bool)
	result := make([]int, 0)

	for _, num := range nums {
		if !seen[num] {
			seen[num] = true
			result = append(result, num)
		}
	}
	return result
}

// Struct and Interface Implementation
type Circle struct {
	Radius float64
}

func (c *Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c *Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

func (c *Circle) Scale(factor float64) {
	c.Radius *= factor
}

// Map Operations
func WordFrequency(text string) map[string]int {
	words := strings.Fields(text)
	frequency := make(map[string]int)

	for _, word := range words {
		frequency[word]++
	}
	return frequency
}

// Error Handling
func GetElement(arr []int, index int) (int, error) {
	if index < 0 || index >= len(arr) {
		return 0, errors.New("index out of bounds")
	}
	return arr[index], nil
}

// Add these implementations
type ThreadSafeCounter struct {
	value int
	mu    sync.RWMutex
}

func NewThreadSafeCounter() *ThreadSafeCounter {
	return &ThreadSafeCounter{
		value: 0,
	}
}

func (c *ThreadSafeCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *ThreadSafeCounter) GetValue() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.value
}

func (c *ThreadSafeCounter) SetValue(v int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value = v
}
