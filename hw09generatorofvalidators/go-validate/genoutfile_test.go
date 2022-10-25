package main

import (
	"bytes"
	"log"
	"os"
	"testing"
)

func TestGenOutFile(t *testing.T) {
	const outFileName = "out.go"
	const wantFileName = "testdata/out.go"

	var structs = validStructs{
		"models",
		[]validStruct{
			{"User",
				[]field{
					{"ID", "string", map[string]string{"len": "36"}},
					{"Age", "int", map[string]string{"min": "18", "max": "50"}},
					{"Email", "string", map[string]string{"regexp": `^\\w+@\\w+\\.\\w+$`}},
					{"Role", "string", map[string]string{"in": "admin,stuff"}},
				},
			},
			{"App", []field{
				{"Version", "string", map[string]string{"len": "5"}},
			},
			},
			{"Response", []field{
				{"Code", "int", map[string]string{"in": "200,404,500"}},
			},
			},
		},
	}

	if err := genOutFile(outFileName, structs); err != nil {
		t.Errorf("genOutFile = %v; want = %v", err, nil)
	}
	if ok, _ := cmpFiles(outFileName, wantFileName); !ok {
		t.Errorf("files: %s, %s - not equel", outFileName, wantFileName)
	}
	if err := os.Remove(outFileName); err != nil {
		log.Fatal(err)
	}
}

func cmpFiles(filePath1, filePath2 string) (bool, error) {
	infFile1, err := os.Stat(filePath1)
	if err != nil {
		return false, err
	}

	infFile2, err := os.Stat(filePath2)
	if err != nil {
		return false, err
	}

	if infFile1.Size() != infFile2.Size() {
		return false, nil
	}

	file1, err := os.ReadFile(filePath1)
	if err != nil {
		return false, err
	}

	file2, err := os.ReadFile(filePath2)
	if err != nil {
		return false, err
	}

	return bytes.Equal(file1, file2), nil
}
