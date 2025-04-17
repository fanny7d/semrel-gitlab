package render

import (
	"testing"
)

func TestBump(t *testing.T) {
	tag, err := BumpMessage("t", "tag is {{tag}}")
	if err != nil {
		t.Error(err)
	}
	if tag != "tag is t" {
		t.Fail()
	}
}
