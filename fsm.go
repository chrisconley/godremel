package main

import "fmt"

type Field struct {
    path string
    maxRepetitionLevel int
}

type FsmState struct {
  field Field
  repetitionLevel int
}

// common repetition level of field and barrier
func commonRepetitionLevel(field Field, barrier Field) int {
  if field.maxRepetitionLevel > barrier.maxRepetitionLevel {
    return barrier.maxRepetitionLevel
  } else {
    return field.maxRepetitionLevel
  }
}

func construct(fields [2]Field) map[FsmState]Field {
    var transitions map[FsmState]Field

    for fieldIndex, field := range fields {
        fmt.Printf("%s\n", field)
        maxLevel := field.maxRepetitionLevel

        // Set the barrier to the next field
        // if there there is still a next field
        numFields := len(fields)
        barrier := Field{"end", 0}
        if fieldIndex+1 < numFields {
          barrier = fields[fieldIndex+1]
        }

        barrierLevel := commonRepetitionLevel(field, barrier)


        // for each preField before field whose repetition level is larger than barrierLevel:
        // Walk each prefield starting with the most recent field
        for preFieldIndex := fieldIndex-1; preFieldIndex >= 0; preFieldIndex-- {
            preField := fields[preFieldIndex]
            if preField.maxRepetitionLevel > barrierLevel {
              backLevel := commonRepetitionLevel(field, preField)
              state := FsmState{field, backLevel}
              transitions[state] = preField
            }
            fmt.Printf("preFieldIndex: %d\n", preFieldIndex)
        }

        // for each level in [barrierLevel+1..maxLevel] that lacks transition from field:


        // for each level in [0..barrierLevel]:

        fmt.Printf("%d | %s | %d\n", maxLevel, barrier, barrierLevel)
    }

    return transitions
}

func main() {
    var fields [2]Field
    fields[0] = Field{"id", 0}
    fields[1] = Field{"names.languages.country", 2}

    m := construct(fields)

    fmt.Printf("fsm: %v\n", m)
    fmt.Printf("hello, world\n")
}
