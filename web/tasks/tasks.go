package tasks

import "sync"

type Tasks struct {
	errors     []TaskError
	errorsLock sync.Mutex
}

func (t *Tasks) Append(name string, f func() error) {
	go func() {
		err := f()
		t.errorsLock.Lock()
		defer t.errorsLock.Unlock()
		t.errors = append(t.errors, TaskError{
			Name:  name,
			error: err,
		})
	}()
}

type TaskError struct {
	Name string
	error
}

func (t *Tasks) Errors() []TaskError {
	return t.errors
}
