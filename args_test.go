package golang_args

import (
	"strings"
	"testing"
)

func getCallback(field *bool) func(...string) {
	return func(s ...string) {
		*field = true
	}
}

func getTestingArgs() []string {
	str := "--name MyName --files a b c d -sptr 12 12 12"
	return strings.Split(str, " ")
}

func getTestingArgMap() map[string][]string {
	ret := make(map[string][]string)
	ret["name"] = []string{"MyName"}
	ret["files"] = []string{"a", "b", "c", "d"}
	ret["s"] = []string{}
	ret["p"] = []string{}
	ret["t"] = []string{}
	ret["r"] = []string{"12", "12", "12"}
	return ret
}

func getParser(callback []func(...string)) (*ArgsParser, error) {
	builder := InitParserBuilder()
	builder.AddElement(callback[0], func(u uint) bool { return u > 0 }, "name", true)
	builder.AddElement(callback[1], func(u uint) bool { return u > 2 }, "files", true)
	builder.AddElement(callback[2], func(u uint) bool { return true }, "sTest", false, "s")
	builder.AddElement(callback[3], func(u uint) bool { return true }, "pTest", false, "p")
	builder.AddElement(callback[4], func(u uint) bool { return true }, "tTest", false, "t")
	builder.AddElement(callback[5], func(u uint) bool { return u > 1 }, "rTest", false, "r")
	return builder.Construct()
}

func getTestCase(t *testing.T, argMap map[string][]string, callbackIfErr func(error)) []bool {
	var checkList [6]bool
	var callbacks [6]func(...string)
	for i := 0; i < 6; i++ {
		checkList[i] = false
		callbacks[i] = getCallback(&checkList[i])
	}
	parser, err := getParser(callbacks[:])
	if err != nil {
		t.Errorf("error of correct parser building %s", err.Error())
	}
	err = parser.Parse(argMap)
	if err != nil {
		callbackIfErr(err)
	}
	return checkList[:]
}

func TestParserCorrectStr(t *testing.T) {
	var checkList [6]bool
	var callbacks [6]func(...string)
	for i := 0; i < 6; i++ {
		checkList[i] = false
		callbacks[i] = getCallback(&checkList[i])
	}
	parser, err := getParser(callbacks[:])
	if err != nil {
		t.Errorf("error of correct parser building %s", err.Error())
	}
	err = parser.ParseArgs(getTestingArgs()...)
	if err != nil {
		t.Errorf("error of parsing %s", err.Error())
	}
	for i, point := range checkList {
		if !point {
			t.Errorf("callback №%d wasn't executed", i)
		}
	}
}

func TestParserCorrectMap(t *testing.T) {
	ifErr := func(err error) { t.Errorf("error of parsing %s", err.Error()) }
	checkList := getTestCase(t, getTestingArgMap(), ifErr)
	for i, point := range checkList {
		if !point {
			t.Errorf("callback №%d wasn't executed", i)
		}
	}
}

func TestParserWithoutUnecessaryArgs(t *testing.T) {
	ifErr := func(err error) { t.Errorf("error of parsing %s", err.Error()) }
	argMap := getTestingArgMap()
	delete(argMap, "p")
	checkList := getTestCase(t, argMap, ifErr)
	for i, point := range checkList {
		if i == 3 {
			if point {
				t.Errorf("callback №%d was executed, but shouldn't to", i)
			}
			continue
		}
		if !point {
			t.Errorf("callback №%d wasn't executed", i)
		}
	}
}

func TestParserWithUnknownArgs(t *testing.T) {
	var wasError bool
	ifErr := func(err error) { wasError = true }
	argMap := getTestingArgMap()
	argMap["shit"] = []string{"shit"}
	getTestCase(t, argMap, ifErr)
	if !wasError {
		t.Errorf("unknown arg pass")
	}
}

func TestParserWithoutNeccesaryArgs(t *testing.T) {
	var wasError bool
	ifErr := func(err error) { wasError = true }
	argMap := getTestingArgMap()
	delete(argMap, "files")
	chekList := getTestCase(t, argMap, ifErr)
	if !wasError {
		t.Errorf("missing of necessary arg pases")
	}
	for i, point := range chekList {
		if point {
			t.Errorf("called callback №%d, but it shouldn't to", i)
		}
	}
}

func TestParserWithIncorrectNumOfArgs(t *testing.T) {
	var wasError bool
	ifErr := func(err error) { wasError = true }
	argMap := getTestingArgMap()
	argMap["files"] = []string{}
	getTestCase(t, argMap, ifErr)
	if !wasError {
		t.Errorf("incorrect num of arg passes")
	}
}
