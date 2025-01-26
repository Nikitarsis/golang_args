package golang_args

import (
	"errors"
	"fmt"
	"strings"
)

// Simple function that divides args to map that associates keys and additional argument
func DivideArgs(args ...string) (map[string][]string, error) {
	result := make(map[string][]string)
	var currentKey string

	//Checks every argument
	for _, arg := range args {
		var check bool
		//If argument has duble minus, argument is count as single key
		if arg, check = strings.CutPrefix(arg, "--"); check {
			currentKey = arg
			//Checking incorrect key
			if currentKey == "" {
				return nil, errors.New("invalid argument key")
			}
			_, exists := result[currentKey]
			if exists {
				return nil, fmt.Errorf("duplicate argument key: %s", currentKey)
			}
			result[currentKey] = []string{}
			continue
		}
		//If argument has single minus, every letter interpretes as key
		if arg, check = strings.CutPrefix(arg, "-"); check {
			currentKey = arg
			keys := strings.Split(currentKey, "")
			if len(keys) == 0 {
				return nil, errors.New("invalid argument key")
			}
			for _, currentKey = range keys {
				//Checking incorrect keys
				if currentKey == "" {
					return nil, errors.New("invalid argument key")
				}

				if _, exists := result[currentKey]; exists {
					return nil, fmt.Errorf("duplicate argument key: %s", currentKey)
				}
				result[currentKey] = []string{}
			}
			continue
		}
		if currentKey == "" {
			return nil, fmt.Errorf("value provided without a key: %s", arg)
		}
		result[currentKey] = append(result[currentKey], arg)
	}
	return result, nil
}
