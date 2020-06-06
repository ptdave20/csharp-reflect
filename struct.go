package csharp_reflect

import (
	"fmt"
	"io/ioutil"
	"path"
	"reflect"
	"strings"
)

type Options struct {
	Namespace           string
	JsonProperty        bool
	OverrideListToArray bool
	types               map[reflect.Kind]string
	BaseUsings          []string
	IdentType           int
	IdentSpacing        int
	OutputPath          string
	SingleFile          bool
}

const (
	IDENT_TAB   = 0
	IDENT_SPACE = 1
)

func New(namespace string) *Options {
	opt := &Options{
		types:               make(map[reflect.Kind]string),
		Namespace:           namespace,
		JsonProperty:        false,
		OverrideListToArray: false,
		IdentType:           IDENT_TAB,
		IdentSpacing:        4,
		SingleFile:          false,
	}

	opt.BaseUsings = []string{
		"using System;",
		"using System.Collection.Generics",
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

	return opt
}

func FieldToCSharpType(t reflect.Type, f reflect.StructField, processed map[string]string) (string, bool) {
	retFmt := "%s"
	vType := ""
	unknown := true
	array := false
	if strings.Index(f.Type.String(), "[]") > -1 {
		// We have an array
		retFmt = "List<%s>"
		array = true
	}

	tName := f.Type.Name()
	if array {
		tName = f.Type.Elem().Name()
	}
	vType = tName

	switch tName {
	case "int":
		vType = "int"
		unknown = false
		break
	case "int32":
		vType = "int"
		unknown = false
		break
	case "int16":
		vType = "Int16"
		unknown = false
		break
	case "uint16":
		vType = "UInt16"
		unknown = false
		break
	case "uint":
		vType = "unsigned int"
		unknown = false
		break
	case "uint32":
		vType = "unsigned int"
		unknown = false
		break
	case "int64":
		vType = "long"
		unknown = false
		break
	case "uint64":
		vType = "long"
		unknown = false
		break
	case "bool":
		vType = "bool"
		unknown = false
		break
	case "byte":
		vType = "byte"
		unknown = false
		break
	case "int8":
		vType = "byte"
		unknown = false
		break
	case "string":
		vType = "string"
		unknown = false
		break
	case "float32":
		vType = "float"
		unknown = false
		break
	}

	if unknown {
		// Is this object in our processed list?
		for k := range processed {
			if k == tName {
				unknown = false
				vType = tName
				break
			}
		}
	}

	return fmt.Sprintf(retFmt, vType), unknown
}

func processType(t reflect.Type, namespace string, jsonProperty bool, processed map[string]string) map[string]string {
	typeName := t.Name()
	if strings.Contains(typeName, ".") {
		typeName = strings.Split(typeName, ".")[1]
	}
	classStrFmt := "\tpublic class " + typeName + " {\r\n%s\t}"
	properties := ""

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fieldType, unknown := FieldToCSharpType(t, f, processed)

		if unknown {
			processed = processType(f.Type, namespace, jsonProperty, processed)
		}

		fieldFmt := "\t\tpublic %s %s { get; set; }\r\n"
		properties += fmt.Sprintf(fieldFmt, fieldType, f.Name)
	}

	class := fmt.Sprintf(classStrFmt, properties)

	processed[typeName] = class

	return processed
}

func processObject(object interface{}, namespace string, jsonProperty bool, processed map[string]string) map[string]string {
	t := reflect.TypeOf(object)
	typeName := t.String()
	if strings.Contains(typeName, ".") {
		typeName = strings.Split(typeName, ".")[1]
	}
	classStrFmt := "\tpublic class " + typeName + " {\r\n%s\t}"
	properties := ""

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fieldType, unknown := FieldToCSharpType(t, f, processed)

		if unknown {
			processed = processType(f.Type, namespace, jsonProperty, processed)
		}

		fieldFmt := "\t\tpublic %s %s { get; set; }\r\n"
		properties += fmt.Sprintf(fieldFmt, strings.TrimSpace(fieldType), f.Name)
	}

	class := fmt.Sprintf(classStrFmt, properties)

	processed[typeName] = class

	return processed
}

func OutputToCSharp(obj interface{}, namespace string, jsonProperty bool, outputPath string) error {
	var typeObject map[string]string = make(map[string]string)
	typeObject = processObject(obj, namespace, jsonProperty, typeObject)
	// Add usings
	for k, v := range typeObject {
		usings := "using System;\r\n"
		// Check for collection
		if strings.Contains(v, "List<") {
			usings += "using System.Collections.Generic;\r\n"
		}
		typeObject[k] = fmt.Sprintf("%s\r\nnamespace %s {\r\n%s\r\n}", usings, namespace, v)
		ioutil.WriteFile(path.Join(outputPath, k+".cs"), []byte(typeObject[k]), 0644)
	}
	return nil
}

func ReadPackage() {

}
