package workflow

import (
	"errors"
	"fmt"
	"time"
)

const (
	// RetryCount is the number of retrys after 502 erros
	RetryCount = 3
	// RetryDelay is the number of seconds to wait before a retry
	RetryDelay = 10
)

// ActionError wraps the original error with a flag
// indicating if the action should be retried
type ActionError struct {
	retry bool
	err   error
}

// Error implements error interface
func (e *ActionError) Error() string {
	return e.err.Error()
}

// NewActionError creates action error
func NewActionError(err error, retry bool) *ActionError {
	return &ActionError{
		retry: retry,
		err:   err,
	}
}

// Action must be idemponent or return an error if that's not possible.
// It is possible that Undo is called before Do.
type Action interface {
	Do() *ActionError
	Undo() error
}

func doActions(actions []Action) (bool, error) {
	for _, action := range actions {
		aErr := action.Do()
		if aErr != nil {
			return aErr.retry, aErr.err
		}
	}
	return false, nil
}

// Apply tries to Do all actions.
// Abort action sequence immediately, if an error occurs
// If error is retryable, retry the sequence RetryCount times,
// with RetryDelay seconds between retries.
// If all retrys fail, rollback to clean up.
func Apply(actions []Action) error {
	var err error
	for i := 0; i <= RetryCount; i++ {
		if i > 0 {
			fmt.Printf("Waiting %d seconds before retry...", RetryDelay)
			time.Sleep(RetryDelay * time.Second)
			fmt.Printf("Starting retry %d/%d\n", i, RetryCount)
		}
		retry, err := doActions(actions)
		if err == nil {
			return nil
		}
		fmt.Println(err.Error())
		if !retry {
			fmt.Println(`
Not trying to continue after this kind of error.
If you think recovery should be attempted, please create an issue
at https://gitlab.com/juhani/go-semrel-gitlab/issues/new
or by email incoming+juhani/go-semrel-gitlab@incoming.gitlab.com`)
			rollbackErr := rollback(actions)
			if rollbackErr != nil {
				fmt.Println(rollbackErr.Error())
			}
			return errors.New("workflow execution failed")
		}
	}
	return err
}

// rollback tries to Undo all actions.
//  list of errors is returned
func rollback(actions []Action) error {
	errorCount := 0
	if len(actions) == 0 {
		return nil
	}
	for i := len(actions) - 1; i >= 0; i-- {
		err := actions[i].Undo()
		if err != nil {
			fmt.Println(err.Error())
			errorCount++
		}
	}
	if errorCount > 0 {
		return fmt.Errorf("%d errors during rollback", errorCount)
	}
	return nil
}
