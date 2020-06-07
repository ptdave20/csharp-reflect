package csharp_reflect

import (
	"fmt"
	tt "github.com/morphar/go-texttools"
	"io/ioutil"
	"log"
	"path"
	"reflect"
	"sort"
	"strings"
)

type Options struct {
	Namespace           string
	JsonProperty        bool
	OverrideListToArray bool
	types               map[reflect.Kind]string
	BaseUsings          []string
	IndentType          INDENT
	IndentSpacing       int
	OutputPath          string
	SingleFile          bool
	Converted           map[reflect.Type]bool
	SortProperties      bool
}

type INDENT int

const (
	INDENT_TAB INDENT = iota
	INDENT_SPACE
)

func New(namespace string) *Options {
	opt := &Options{
		types:               make(map[reflect.Kind]string),
		Namespace:           namespace,
		JsonProperty:        false,
		OverrideListToArray: false,
		IndentType:          INDENT_TAB,
		IndentSpacing:       4,
		SingleFile:          false,
		Converted:           make(map[reflect.Type]bool),
		SortProperties:      true,
	}

	opt.BaseUsings = []string{
		"using System;\r\n",
		"using System.Collection.Generics;\r\n",
	}

	opt.types[reflect.Int] = "int"
	opt.types[reflect.Int8] = "byte"
	opt.types[reflect.Int16] = "Int16"
	opt.types[reflect.Int32] = "int"
	opt.types[reflect.Int64] = "long"
	opt.types[reflect.Uint] = "unsigned int"
	opt.types[reflect.Uint8] = "byte"
	opt.types[reflect.Uint16] = "UInt16"
	opt.types[reflect.Uint32] = "unsigned int"
	opt.types[reflect.Uint64] = "unsigned long"
	opt.types[reflect.Bool] = "bool"
	opt.types[reflect.Float32] = "float"
	opt.types[reflect.Float64] = "double"
	opt.types[reflect.String] = "string"

	return opt
}

func outputEnumerable(objectType string, objectName string, options *Options) string {
	if options.OverrideListToArray {
		return outputProperty(objectType+"[]", objectName, options)
	}
	return outputProperty("List<"+objectType+">", objectName, options)
}

func outputProperty(objectType string, objectName string, options *Options) string {
	if strings.Contains(objectType, ".") {
		objectType = strings.Split(objectType, ".")[1]
	}
	return fmt.Sprintf("public %s %s { get; set; }\r\n", objectType, objectName)
}

func convertField(f reflect.StructField, options *Options) string {
	if f.Type.Kind() == reflect.Ptr {

	}
	if f.Type.Kind() == reflect.Slice {
		// We have a slice of somekind
		return outputEnumerable(f.Type.Elem().Name(), f.Name, options)
	}
	if f.Type.Kind() == reflect.Struct {
		ConvertType(f.Type, options)

		// Have we already processed this?
		if options.Converted[f.Type] {
			return outputProperty(f.Type.String(), f.Name, options)
		}

	}

	propType := options.types[f.Type.Kind()]
	if len(propType) > 0 {
		return outputProperty(propType, f.Name, options)
	}

	return ""
}

func indent(multiple int, value string, options *Options) string {
	if value == "" {
		return ""
	}
	ret := ""
	if options.IndentType == INDENT_TAB {
		for i := 0; i < multiple; i++ {
			ret += "\t"
		}
	} else if options.IndentType == INDENT_SPACE {
		for i := 0; i < multiple; i++ {
			for s := 0; s < options.IndentSpacing; s++ {
				ret += " "
			}
		}
	}
	return ret + value
}

func ConvertType(t reflect.Type, options *Options) {
	if options.Converted[t] {
		return
	}
	typeName := t.Name()
	options.Converted[t] = true
	if strings.Contains(typeName, ".") {
		typeName = strings.Split(typeName, ".")[1]
	}

	var classProps = make(map[string]string)

	// Pascal Case the object
	typeName = tt.PascalCase(typeName)
	for fieldIndex := 0; fieldIndex < t.NumField(); fieldIndex++ {
		field := t.Field(fieldIndex)

		propertyName := field.Name
		classProps[propertyName] = indent(2, convertField(field, options), options)
	}
	keys := make([]string, 0, len(classProps))
	for k := range classProps {
		keys = append(keys, k)
	}

	// Sort the fields
	if options.SortProperties {
		sort.Strings(keys)
	}

	props := ""
	for _, k := range keys {
		props += classProps[k]
	}

	// Prepare a the class
	classPrefix := "public class " + typeName + " {\r\n"
	classPostfix := "}"
	classPrefix = indent(1, classPrefix, options)
	classPostfix = indent(1, classPostfix, options)
	class := classPrefix + props + classPostfix

	// Attach namespace
	namespace := fmt.Sprintf("namespace %s {\r\n%s\r\n}", options.Namespace, class)

	// Attach usings
	result := ""
	for _, v := range options.BaseUsings {
		result += v
	}
	result += "\r\n" + namespace

	filePath := path.Join(options.OutputPath, typeName+".cs")
	err := ioutil.WriteFile(filePath, []byte(result), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
func ConvertObject(t interface{}, options *Options) {
	rT := reflect.TypeOf(t)
	ConvertType(rT, options)
}
