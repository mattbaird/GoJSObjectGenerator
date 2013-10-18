package jsify

import (
	"os"
	"testing"
)

func Test_NoErrors(t *testing.T) {
	structs := make([]interface{}, 3)
	structs[0] = &Foo{}
	structs[1] = &Bar{}
	structs[2] = &Baz{}
	_, err := GenerateJavascriptToString(structs)
	if err != nil {
		t.Fail()
	}
}

func Test_MakeFile(t *testing.T) {
	fileName := "./structs.js"
	defer os.Remove(fileName)
	structs := make([]interface{}, 3)
	structs[0] = &Foo{}
	structs[1] = &Bar{}
	structs[2] = &Baz{}
	err := GenerateJavascriptToFile(fileName, structs)
	if err != nil {
		t.Fail()
	}
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		t.Fail()
	}
}
