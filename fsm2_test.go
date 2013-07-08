package fsm2

import (
    //. "fsm2"
    "testing"
)

func TestSqrt(t *testing.T) {
  field1 := fsm2.ProcessedField{}
  field2 := fsm2.ProcessedField{}
  result := fsm2.getCommonRepetitionLevel(field1, field2)
	if result != 1 {
		t.Errorf("Expected %v to be %v", result, 1)
	}
}
