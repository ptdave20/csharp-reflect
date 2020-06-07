package csharp_reflect

import (
	"os"
	"path"
	"testing"
)

func FilesExist(outputPath string, t *testing.T, files ...string) bool {
	for _, file := range files {
		filePath := path.Join(outputPath, file)
		if _, err := os.Stat(filePath); err != nil && os.IsNotExist(err) {
			t.Failed()
			return false
		}
	}
	return true
}

type SimpleTest struct {
	A int
	B bool
	C float32
}

func TestSimple(t *testing.T) {
	var testStruct SimpleTest
	options := New("Test")
	options.OutputPath = "./TestOutput/SimpleTest"
	ConvertObject(testStruct, options)
	// Check to see which files are in the directory
	FilesExist(options.OutputPath, t, "SimpleTest.cs")
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
	options.OutputPath = "./TestOutput/EmbedTest"
	ConvertObject(testStruct, options)

	// Check to see which files are in the directory
	FilesExist(options.OutputPath, t, "SimpleTest.cs", "EmbedTest.cs")
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
	options.OutputPath = "./TestOutput/ArrayTest"
	ConvertObject(testStruct, options)

	// Check to see which files are in the directory
	FilesExist(options.OutputPath, t, "SimpleTest.cs", "ArrayTest.cs")
}
