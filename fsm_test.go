package go_dremel

import (
    "fmt"
    "testing" //import go package for testing related functionality
    "io/ioutil"
    "encoding/json"
    )

func assertFsmTransition(fsm FSM2, currentField ProcessedField, rLevel int, destination ProcessedField) (bool, string) {
  if fsm[currentField][rLevel] != destination {
    return false, fmt.Sprintf("Expected (%v, %v) to be %v, but it was %v\n", currentField, rLevel, destination, fsm[currentField][rLevel])
  } else {
    return true, ""
  }
}

func TestConstructFSMMimimum(t *testing.T) {
  file, e := ioutil.ReadFile("./docs.json")
  if e != nil {
      t.Errorf("File error: %v\n", e)
  }

  var schema Field
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
  if exists, err := assertFsmTransition(fsm, idField, 0, countryField); !exists {
    t.Errorf(err)
  }

  //countryField
  if exists, err := assertFsmTransition(fsm, countryField, 1, countryField); !exists {
    t.Errorf(err)
  }
  if exists, err := assertFsmTransition(fsm, countryField, 2, countryField); !exists {
    t.Errorf(err)
  }
  if exists, err := assertFsmTransition(fsm, countryField, 0, EndField); !exists {
    t.Errorf(err)
  }
}

func TestConstructFSMPartial(t *testing.T) {
  file, e := ioutil.ReadFile("./docs.json")
  if e != nil {
      t.Errorf("File error: %v\n", e)
  }

  var schema Field
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
  if exists, err := assertFsmTransition(fsm, idField, 0, codeField); !exists {
    t.Errorf(err)
  }

  //codeField
  if exists, err := assertFsmTransition(fsm, codeField, 0, countryField); !exists {
    t.Errorf(err)
  }
  if exists, err := assertFsmTransition(fsm, codeField, 1, countryField); !exists {
    t.Errorf(err)
  }
  if exists, err := assertFsmTransition(fsm, codeField, 2, countryField); !exists {
    t.Errorf(err)
  }

  //countryField
  if exists, err := assertFsmTransition(fsm, countryField, 2, codeField); !exists {
    t.Errorf(err)
  }
  if exists, err := assertFsmTransition(fsm, countryField, 1, countryField); !exists {
    t.Errorf(err)
  }
  if exists, err := assertFsmTransition(fsm, countryField, 0, EndField); !exists {
    t.Errorf(err)
  }
}


func TestConstructFSMFull(t *testing.T) {
  file, e := ioutil.ReadFile("./docs.json")
  if e != nil {
      t.Errorf("File error: %v\n", e)
  }

  var schema Field
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
  if exists, err := assertFsmTransition(fsm, idField, 0, backwardField); !exists {
    t.Errorf(err)
  }

  //backwardField
  if exists, err := assertFsmTransition(fsm, backwardField, 1, backwardField); !exists {
    t.Errorf(err)
  }
  if exists, err := assertFsmTransition(fsm, backwardField, 0, forwardField); !exists {
    t.Errorf(err)
  }

  //forwardField
  if exists, err := assertFsmTransition(fsm, forwardField, 1, forwardField); !exists {
    t.Errorf(err)
  }
  if exists, err := assertFsmTransition(fsm, forwardField, 0, codeField); !exists {
    t.Errorf(err)
  }

  //codeField
  if exists, err := assertFsmTransition(fsm, codeField, 0, countryField); !exists {
    t.Errorf(err)
  }
  if exists, err := assertFsmTransition(fsm, codeField, 1, countryField); !exists {
    t.Errorf(err)
  }
  if exists, err := assertFsmTransition(fsm, codeField, 2, countryField); !exists {
    t.Errorf(err)
  }

  //countryField
  if exists, err := assertFsmTransition(fsm, countryField, 2, codeField); !exists {
    t.Errorf(err)
  }
  if exists, err := assertFsmTransition(fsm, countryField, 0, urlField); !exists {
    t.Errorf(err)
  }
  if exists, err := assertFsmTransition(fsm, countryField, 1, urlField); !exists {
    t.Errorf(err)
  }

  //urlField
  if exists, err := assertFsmTransition(fsm, urlField, 1, codeField); !exists {
    t.Errorf(err)
  }
  if exists, err := assertFsmTransition(fsm, urlField, 0, EndField); !exists {
    t.Errorf(err)
  }
}
