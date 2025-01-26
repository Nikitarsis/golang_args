package golang_args

import (
	"errors"
	"fmt"
)

/*
Parser that analyze arguments and use callback functions func(...string) to change program.
*/
type ArgsParser struct {
	varnames             map[string]string          //Possible names of key including singlekey form. Using for conversion to one key name
	possibleNArgsChecker map[string]func(uint) bool //Function checks number of arguments
	functionMapper       map[string]func(...string) //Callback function that call if arg exists. Argument of function is additional arguments to key
	checklist            map[string]bool            //List of necessary keys
}

// Default parser
func newParserDefault() *ArgsParser {
	return &ArgsParser{make(map[string]string), make(map[string]func(uint) bool), make(map[string]func(...string)), make(map[string]bool)}
}

// Function checks can be arg entity be appened to Parser
func (ap ArgsParser) checkEntity(entity argEntity) (bool, string) {
	//Check every pseudonym
	for _, pseudonym := range entity.Pseudonyms {
		if refer, check := ap.varnames[pseudonym]; check {
			return false, fmt.Sprintf("Pseudonym %s already exists and refer to %s", pseudonym, refer)
		}
	}
	//Check argument name
	if _, check := ap.varnames[entity.ArgName]; check {
		return false, fmt.Sprintf("Keyname %s already exist", entity.ArgName)
	}
	//Check number of args
	if _, check := ap.possibleNArgsChecker[entity.ArgName]; check {
		return false, fmt.Sprintf("Critical error. Keyname %s already refer to check function, but somehow doesn't exist.", entity.ArgName)
	}
	//Check argument name in functions
	if _, check := ap.functionMapper[entity.ArgName]; check {
		return false, fmt.Sprintf("Critical error. Keyname %s already refer to map function, but somehow doesn't exist and don't have check function", entity.ArgName)
	}
	//Check argument name in checklist
	if _, check := ap.checklist[entity.ArgName]; check {
		return false, fmt.Sprintf("Critical error. Keyname %s already contains in checklist, but somehow doesn't exist", entity.ArgName)
	}
	return true, "Ok"
}

// Checks and add new Entity to Parser
func (ap *ArgsParser) addEntity(entity argEntity) error {
	//if argument cannot be insert into Parser, error returns
	if check, msg := ap.checkEntity(entity); !check {
		return errors.New("Entity collision due to: " + msg)
	}

	ap.possibleNArgsChecker[entity.ArgName] = entity.NargsChecker
	ap.functionMapper[entity.ArgName] = entity.AssociatedFunction
	ap.varnames[entity.ArgName] = entity.ArgName
	for _, element := range entity.Pseudonyms {
		ap.varnames[element] = entity.ArgName
	}
	ap.checklist[entity.ArgName] = entity.IsNecessary

	return nil
}

// Parse arguments in form map[keymap string]arguments []string
// Function is public due to using Parser with others syntax analyzers
func (ap ArgsParser) Parse(argsMap map[string][]string) error {
	for key, isNecessary := range ap.checklist {
		_, check := argsMap[key]
		if isNecessary && !check {
			return fmt.Errorf("should be key %s", key)
		}
	}
	for key, values := range argsMap {
		pseudonym, check := ap.varnames[key]
		if !check {
			return fmt.Errorf("key %s doesn't exist", key)
		}
		checker, check := ap.possibleNArgsChecker[pseudonym]
		if !check {
			return fmt.Errorf("no checker for key %s", pseudonym)
		}
		function, check := ap.functionMapper[pseudonym]
		if !check {
			return fmt.Errorf("no such function with name %s", pseudonym)
		}
		size := uint(len(values))
		if !checker(size) {
			return fmt.Errorf("incorrect number of arguments %s", key)
		}
		function(values...)
	}
	return nil
}

// Parse arguments and return error if something goes wrong
func (ap ArgsParser) ParseArgs(args ...string) error {
	argsMap, error := DivideArgs(args...)
	if error != nil {
		return error
	}
	return ap.Parse(argsMap)
}

// Constuct Parser from entities
func constructParser(entities ...argEntity) (*ArgsParser, error) {
	var ret = newParserDefault()
	for _, entity := range entities {
		err := ret.addEntity(entity)
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}
