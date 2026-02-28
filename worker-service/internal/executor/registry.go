package executor

import (
	"fmt"
	"sync"
)

type Registry interface {
	Register(jobType string, exec Executor) error
	Get(jobType string) (Executor, error)
	List() []string
}

type registry struct {
	mutex     sync.RWMutex
	executors map[string]Executor
}

func NewRegistry() Registry {
	return &registry{
		executors: make(map[string]Executor),
	}
}

func (r *registry) Register(jobType string, exec Executor) error {
	if jobType == "" {
		return fmt.Errorf("jobType cannot be empty")
	}

	if exec == nil {
		return fmt.Errorf("executor cannot be nil")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.executors[jobType]; exists {
		return fmt.Errorf("executor already registered for job type: %s", jobType)
	}

	r.executors[jobType] = exec
	return nil
}

func (r *registry) Get(jobType string) (Executor, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	exec, ok := r.executors[jobType]
	if !ok {
		return nil, fmt.Errorf("no executor registered for job type: %s", jobType)
	}

	return exec, nil
}

func (r *registry) List() []string {

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	keys := make([]string, 0, len(r.executors))

	for k := range r.executors {
		keys = append(keys, k)
	}

	return keys
}
