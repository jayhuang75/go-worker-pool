package worker

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWorkerPool(t *testing.T) {
	pool := NewWorkerPool(10)
	assert.False(t, pool.IsCompleted())
}

func TestLog(t *testing.T) {
	pool := NewWorkerPool(10)
	pool.Log(false)
	pool.logging = false
}

func ResourceProcessor(resource interface{}) error {
	fmt.Printf("Resource processor got: %s", resource)
	fmt.Println()
	return nil
}

func ResultProcessor(result Result) error {
	fmt.Printf("Result processor got: %s", result.Err)
	fmt.Println()
	return nil
}

func TestPool_Start(t *testing.T) {
	strings := []string{"first", "second", "first", "second", "first", "second", "first", "second", "first", "second"}
	resources := make([]interface{}, len(strings))
	for i, s := range strings {
		resources[i] = s
	}

	pool := NewWorkerPool(10)
	pool.Start(resources, ResourceProcessor, ResultProcessor)
}
