package go_dremel

import (
    "fmt"
    "testing" //import go package for testing related functionality
    )

func TestSqrt(t *testing.T) {
  fmt.Printf("hi")
  field1 := ProcessedField{}
  field2 := ProcessedField{}
  result := getCommonRepetitionLevel(field1, field2)
	if result != 2 {
		t.Errorf("Expected %v to be %v", result, 1)
	}
}
