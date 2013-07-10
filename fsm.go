package go_dremel

import "fmt"

type FsmState struct {
  field ProcessedField
  repetitionLevel int
}

func ConstructFSM(fields []ProcessedField) map[FsmState]ProcessedField {
    transitions := map[FsmState]ProcessedField{}

    for fieldIndex, field := range fields {
        fmt.Printf("%s\n", field)
        maxLevel := field.MaxRepetitionLevel()

        // Set the barrier to the next field
        // if there there is still a next field
        numFields := len(fields)
        barrier := ProcessedField{"end", "", "", &ProcessedField{}}
        if fieldIndex+1 < numFields {
          barrier = fields[fieldIndex+1]
        }

        barrierLevel := GetCommonRepetitionLevel(field, barrier)


        // for each preField before field whose repetition level is larger than barrierLevel:
        // Walk each prefield starting with the most recent field
        for preFieldIndex := fieldIndex-1; preFieldIndex >= 0; preFieldIndex-- {
            preField := fields[preFieldIndex]
            if preField.MaxRepetitionLevel() > barrierLevel {
              backLevel := GetCommonRepetitionLevel(field, preField)
              state := FsmState{field, backLevel}
              fmt.Printf("preField %v\n", preField)
              transitions[state] = preField
            }
            fmt.Printf("preFieldIndex: %d\n", preFieldIndex)
            fmt.Printf("preField.MaxRepetitionLevel: %d\n", preField.MaxRepetitionLevel())
            fmt.Printf("barrier level: %d\n", barrierLevel)
        }

        // for each level in [barrierLevel+1..maxLevel] that lacks transition from field:


        // for each level in [0..barrierLevel]:

        fmt.Printf("%d | %s | %d\n", maxLevel, barrier, barrierLevel)
    }

    return transitions
}
