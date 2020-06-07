package csharp_reflect

import "testing"

type SimpleTest struct {
	A int
	B bool
	C float32
}

func TestSimple(t *testing.T) {
	var testStruct SimpleTest
	//OutputToCSharp(testStruct, "test", false, "./TestOutput")
	options := New("Test")
	options.OutputPath = "./TestOutput"
	ConvertObject(testStruct, options)
}

type EmbedTest struct {
	Simple SimpleTest
	A      int
	B      bool
	C      float32
}

func TestEmbed(t *testing.T) {
	var testStruct EmbedTest
	//OutputToCSharp(testStruct, "test", false, "./TestOutput")
	options := New("Test")
	options.OutputPath = "./TestOutput"
	ConvertObject(testStruct, options)
}

type ArrayTest struct {
	A []int
	B bool
	C SimpleTest
	D []SimpleTest
}

func TestArray(t *testing.T) {
	var testStruct ArrayTest
	//OutputToCSharp(testStruct, "test", false, "./TestOutput")
	options := New("Test")
	options.OutputPath = "./TestOutput"
	ConvertObject(testStruct, options)
}
