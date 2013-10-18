package jsify

import (
	"bytes"
	"errors"
	"io/ioutil"
	"reflect"
	"strings"
	"text/template"
)

func GenerateJavascriptToFile(file string, structs []interface{}) error {
	source, err := GenerateJavascriptToString(structs)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, []byte(source), 0644)
}

func GenerateJavascriptToString(structs []interface{}) (string, error) {
	structified := make([]interface{}, len(structs))
	for i, mystruct := range structs {
		s, err := structify(mystruct)
		if err != nil {
			return "", err
		}
		structified[i] = s
	}
	// Create a new template and parse the jsClass into it.
	t := template.Must(template.New("jsClass").Funcs(funcMap).Parse(jsClass))
	source := ""
	// Execute the template for each struct.
	for _, r := range structified {
		var sourceBytes bytes.Buffer
		err := t.Execute(&sourceBytes, r)
		if err != nil {
			return "", err
		}
		source += sourceBytes.String()
	}
	return source, nil
}

type Struct struct {
	ClassName string
	Fields    []Field
}

type Field struct {
	Name           string
	JavascriptType string
}

var funcMap = template.FuncMap{
	"variables": func(fields []Field) string {
		names := []string{}
		for _, field := range fields {
			names = append(names, camelCase(field.Name))
		}
		return strings.Join(names, ",")
	},
	"tolower": func(s string) string {
		return strings.ToLower(s)
	},
	"camelCase": func(s string) string {
		return camelCase(s)
	},
}

func camelCase(input string) string {
	if len(input) == 0 {
		return input
	} else if len(input) == 1 {
		return strings.ToLower(input)
	} else {
		return strings.ToLower(input[0:1]) + input[1:len(input)]
	}
}

func structify(obj interface{}) (Struct, error) {
	if obj != nil {
		name := reflect.TypeOf(obj).String()
		index := strings.LastIndex(name, ".")
		if index > 0 {
			name = name[index+1 : len(name)]
		}
		val := reflect.ValueOf(obj).Elem()
		fields := []Field{}
		for i := 0; i < val.NumField(); i++ {
			// valueField := val.Field(i)
			// tag := typeField.Tag
			typeField := val.Type().Field(i)
			typeOfField := typeField.Type.String()
			fields = append(fields, Field{Name: typeField.Name, JavascriptType: getJavascriptType(typeOfField)})
		}
		return Struct{ClassName: name, Fields: fields}, nil
	} else {
		return Struct{}, errors.New("The obj you passed to structify was nil.")
	}
}

func getJavascriptType(from string) string {
	switch from {
	case "int", "int8", "int16", "int32", "int64":
		return "number"
	case "string":
		return "string"
	default:
		return "object"
	}
}

// Define a template.
var jsClass = `

var {{camelCase .ClassName}} = function({{ variables .Fields}}) {

    this.properties = {};
    var numArguments = {{len .Fields}};

    if (arguments.length !== numArguments) {
        throw("Incorrect number of arguments, expected " + numArguments + " to be passed, but received " + arguments.length);
    } 
{{range .Fields}}
    this.__defineSetter__("{{camelCase .Name}}", function({{camelCase .Name}}) {
        if(typeof {{camelCase .Name}} !== '{{.JavascriptType}}') {
            throw("argument '{{camelCase .Name}}' should have been '{{.JavascriptType}}', but instead was " + typeof {{camelCase .Name}});
        } else {
            this.properties.{{camelCase .Name}} = {{camelCase .Name}};
        }
    });
{{end}}
    this.__defineGetter__("json", function() {
        return JSON.stringify(this.properties);
    });

{{range .Fields}}    this.{{camelCase .Name}} = {{camelCase .Name}};
{{end}}
};

`

/**
// Instantiating
var something = new className(myGoo, myFoo, myBar, myBaz);
**/
