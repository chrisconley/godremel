package go_dremel

import (
    "fmt"
    "testing" //import go package for testing related functionality
    "io/ioutil"
    "encoding/json"
    )

func assertFsmTransition(fsm FSM, state FsmState, destination ProcessedField) (bool, string) {
  if fsm[state] != destination {
    return false, fmt.Sprintf("Expected %v to be %v, but it was %v\n", state, destination, fsm[state])
  } else {
    return true, ""
  }
}

func TestConstructFSMMimimum(t *testing.T) {
  file, e := ioutil.ReadFile("./docs.json")
  if e != nil {
      t.Errorf("File error: %v\n", e)
  }

  var schema Schema
  err := json.Unmarshal(file, &schema)
  if err != nil {
    t.Errorf("Json error: %v\n", err)
  }
  pFields := processFields(schema.Fields, []ProcessedField{}, ProcessedField{})
  fields := findFields(pFields, "id", "names.languages.country")
  idField := fields[0]
  countryField := fields[1]
  fsm := ConstructFSM(fields)

  //idField
  if fsm[FsmState{idField, 0}] != countryField {
    t.Errorf("hi")
  }

  //countryField
  if fsm[FsmState{countryField, 1}] != countryField {
    t.Errorf("hi")
  }
  if fsm[FsmState{countryField, 2}] != countryField {
    t.Errorf("hi")
  }
  if fsm[FsmState{countryField, 0}] != EndField {
    t.Errorf("hi")
  }
}

func TestConstructFSMPartial(t *testing.T) {
  file, e := ioutil.ReadFile("./docs.json")
  if e != nil {
      t.Errorf("File error: %v\n", e)
  }

  var schema Schema
  err := json.Unmarshal(file, &schema)
  if err != nil {
    t.Errorf("Json error: %v\n", err)
  }
  pFields := processFields(schema.Fields, []ProcessedField{}, ProcessedField{})
  fields := findFields(pFields, "id", "names.languages.code", "names.languages.country")
  idField := fields[0]
  codeField := fields[1]
  countryField := fields[2]
  fsm := ConstructFSM(fields)

  //idField
  if fsm[FsmState{idField, 0}] != codeField {
    t.Errorf("hi")
  }

  //codeField
  if fsm[FsmState{codeField, 0}] != countryField {
    t.Errorf("hi")
  }
  if fsm[FsmState{codeField, 1}] != countryField {
    t.Errorf("hi")
  }
  if fsm[FsmState{codeField, 2}] != countryField {
    t.Errorf("hi")
  }

  //countryField
  if fsm[FsmState{countryField, 2}] != codeField {
    t.Errorf("hi")
  }
  if fsm[FsmState{countryField, 1}] != countryField {
    t.Errorf("hi")
  }
  if fsm[FsmState{countryField, 0}] != EndField {
    t.Errorf("hi")
  }
}


func TestConstructFSMFull(t *testing.T) {
  file, e := ioutil.ReadFile("./docs.json")
  if e != nil {
      t.Errorf("File error: %v\n", e)
  }

  var schema Schema
  err := json.Unmarshal(file, &schema)
  if err != nil {
    t.Errorf("Json error: %v\n", err)
  }
  pFields := processFields(schema.Fields, []ProcessedField{}, ProcessedField{})
  fields := findFields(pFields, "id", "links.backward", "links.forward",
      "names.languages.code", "names.languages.country", "names.url")
  idField := fields[0]
  backwardField := fields[1]
  forwardField := fields[2]
  codeField := fields[3]
  countryField := fields[4]
  urlField := fields[5]
  fsm := ConstructFSM(fields)

  //idField
  if fsm[FsmState{idField, 0}] != backwardField {
    t.Errorf("hi")
  }

  //backwardField
  if fsm[FsmState{backwardField, 1}] != backwardField {
    t.Errorf("hi")
  }
  if exists, err := assertFsmTransition(fsm, FsmState{backwardField, 0}, forwardField); !exists {
    t.Errorf(err)
  }

  //forwardField
  if fsm[FsmState{forwardField, 1}] != forwardField {
    t.Errorf("hi")
  }
  if exists, err := assertFsmTransition(fsm, FsmState{forwardField, 0}, codeField); !exists {
    t.Errorf(err)
  }

  //codeField
  if exists, err := assertFsmTransition(fsm, FsmState{codeField, 0}, countryField); !exists {
    t.Errorf(err)
  }
  if exists, err := assertFsmTransition(fsm, FsmState{codeField, 1}, countryField); !exists {
    t.Errorf(err)
  }
  if exists, err := assertFsmTransition(fsm, FsmState{codeField, 2}, countryField); !exists {
    t.Errorf(err)
  }

  //countryField
  if exists, err := assertFsmTransition(fsm, FsmState{countryField, 2}, codeField); !exists {
    t.Errorf(err)
  }
  if exists, err := assertFsmTransition(fsm, FsmState{countryField, 0}, urlField); !exists {
    t.Errorf(err)
  }
  if exists, err := assertFsmTransition(fsm, FsmState{countryField, 1}, urlField); !exists {
    t.Errorf(err)
  }

  //urlField
  if exists, err := assertFsmTransition(fsm, FsmState{urlField, 1}, codeField); !exists {
    t.Errorf(err)
  }
  if exists, err := assertFsmTransition(fsm, FsmState{urlField, 0}, EndField); !exists {
    t.Errorf(err)
  }
}
