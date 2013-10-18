package jsify

import (
	"testing"
)

func Test_Example(t *testing.T) {
	structs := make([]interface{}, 3)
	structs[0] = &Foo{}
	structs[1] = &Bar{}
	structs[2] = &Baz{}
	_, err := GenerateJavascriptToString(structs)
	if err != nil {
		t.Fail()
	}
}
