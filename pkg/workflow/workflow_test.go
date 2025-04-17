package workflow

import "testing"

type simple struct {
	done bool
}

func (a *simple) Do() *ActionError {
	a.done = true
	return nil
}

func (a *simple) Undo() error {
	a.done = false
	return nil
}

func TestApplyAndRollback(t *testing.T) {
	var simples []*simple
	var actions []Action
	for _, a := range make([]simple, 5) {
		actions = append(actions, &a)
		simples = append(simples, &a)
	}
	aErr := Apply(actions)
	if aErr != nil {
		t.Fatal(aErr)
	}
	for _, a := range simples {
		if !a.done {
			t.Error("Action not applied")
		}
	}
}
