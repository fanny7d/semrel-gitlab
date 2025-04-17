package workflow

import (
	"errors"
	"testing"
)

// TestAction 是一个测试用的 Action 实现
type TestAction struct {
	doError    *ActionError
	undoError  error
	doCalled   bool
	undoCalled bool
}

func (a *TestAction) Do() *ActionError {
	a.doCalled = true
	return a.doError
}

func (a *TestAction) Undo() error {
	a.undoCalled = true
	return a.undoError
}

func TestApply_Success(t *testing.T) {
	action := &TestAction{}
	err := Apply([]Action{action})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !action.doCalled {
		t.Error("Expected Do to be called")
	}
	if action.undoCalled {
		t.Error("Expected Undo not to be called")
	}
}

func TestApply_Error(t *testing.T) {
	expectedErr := errors.New("test error")
	action := &TestAction{
		doError: NewActionError(expectedErr, false),
	}
	err := Apply([]Action{action})
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if !action.doCalled {
		t.Error("Expected Do to be called")
	}
	if !action.undoCalled {
		t.Error("Expected Undo to be called")
	}
}

func TestApply_RetryableError(t *testing.T) {
	expectedErr := errors.New("test error")
	action := &TestAction{
		doError: NewActionError(expectedErr, true),
	}
	err := Apply([]Action{action})
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if !action.doCalled {
		t.Error("Expected Do to be called")
	}
	if action.undoCalled {
		t.Error("Expected Undo not to be called")
	}
}

func TestRollback_Success(t *testing.T) {
	action := &TestAction{}
	err := rollback([]Action{action})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !action.undoCalled {
		t.Error("Expected Undo to be called")
	}
}

func TestRollback_Error(t *testing.T) {
	expectedErr := errors.New("test error")
	action := &TestAction{
		undoError: expectedErr,
	}
	err := rollback([]Action{action})
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if !action.undoCalled {
		t.Error("Expected Undo to be called")
	}
}
