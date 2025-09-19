package tester

import (
	"context"
	"fmt"
	"sync"
)

// MemoryStore implements Store interface using in-memory storage
type MemoryStore struct {
	testSuites map[string]*TestSuite
	executions map[string]*TestExecution
	mutex      sync.RWMutex
}

// NewMemoryStore creates a new in-memory store
func NewMemoryStore() Store {
	return &MemoryStore{
		testSuites: make(map[string]*TestSuite),
		executions: make(map[string]*TestExecution),
	}
}

// SaveTestSuite saves a test suite
func (ms *MemoryStore) SaveTestSuite(ctx context.Context, suite *TestSuite) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	ms.testSuites[suite.ID] = suite
	return nil
}

// GetTestSuite retrieves a test suite by ID
func (ms *MemoryStore) GetTestSuite(ctx context.Context, suiteID string) (*TestSuite, error) {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	suite, exists := ms.testSuites[suiteID]
	if !exists {
		return nil, fmt.Errorf("test suite with ID %s not found", suiteID)
	}

	return suite, nil
}

// SaveTestExecution saves a test execution
func (ms *MemoryStore) SaveTestExecution(ctx context.Context, execution *TestExecution) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	ms.executions[execution.ID] = execution
	return nil
}

// GetTestExecution retrieves a test execution by ID
func (ms *MemoryStore) GetTestExecution(ctx context.Context, executionID string) (*TestExecution, error) {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	execution, exists := ms.executions[executionID]
	if !exists {
		return nil, fmt.Errorf("test execution with ID %s not found", executionID)
	}

	return execution, nil
}

// UpdateTestExecution updates a test execution
func (ms *MemoryStore) UpdateTestExecution(ctx context.Context, execution *TestExecution) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	if _, exists := ms.executions[execution.ID]; !exists {
		return fmt.Errorf("test execution with ID %s not found", execution.ID)
	}

	ms.executions[execution.ID] = execution
	return nil
}
