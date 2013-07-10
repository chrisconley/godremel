package go_dremel

type FsmState struct {
  field ProcessedField
  repetitionLevel int
}

type FSM map[FsmState]ProcessedField

var EndField = ProcessedField{"end", "", "", &ProcessedField{}}

func ConstructFSM(fields []ProcessedField) FSM {
    fsm := FSM{}

    for fieldIndex, field := range fields {
        maxLevel := field.MaxRepetitionLevel()

        // Set the barrier to the next field
        // if there there is still a next field
        barrier := EndField
        if fieldIndex+1 < len(fields) {
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
              fsm[state] = preField
            }
        }

        // for each level in [barrierLevel+1..maxLevel] that lacks transition from field:
        for level := barrierLevel+1; level <= maxLevel; level++ {
          if _, present := fsm[FsmState{field, level}]; !present {
            //fsm[FsmState{field, level}] = fsm[FsmState{field, level - 1}]
            fsm[FsmState{field, level}] = field // The whitepaper says get field from level-1, but this works
          }
        }

        // for each level in [0..barrierLevel], move to barrier (next field)
        for level := 0; level <= barrierLevel; level++ {
          fsm[FsmState{field, level}] = barrier
        }
    }

    return fsm
}
