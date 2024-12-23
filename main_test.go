package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

type MockDataGenerator struct {
	ShouldFail bool
}

func (d MockDataGenerator) generateCsvData(rows int, fields string, outputDir string, filename string, fileHandler FileHandler, csvWriter FileWriter) error {
	if d.ShouldFail {
		return fmt.Errorf("generateCsvData failed")
	}

	return nil
}

type MockFileHandler struct {
	ShouldFailMkDirAll bool
	ShouldFailCreate   bool
}

func (f MockFileHandler) MkDirAll(path string, perm os.FileMode) error {
	if f.ShouldFailMkDirAll {
		return fmt.Errorf("MkDirAll failed")
	}

	return nil
}

func (f MockFileHandler) Create(name string) (*os.File, error) {
	if f.ShouldFailCreate {
		return nil, fmt.Errorf("Create failed")
	}

	return nil, nil
}

type MockFileWriter struct {
	ShouldFail bool
}

func (w MockFileWriter) Write(row []string, writer *csv.Writer) error {
	if w.ShouldFail {
		return fmt.Errorf("Write failed")
	}

	return nil
}

func TestMain_ErrorCases(t *testing.T) {
	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	tests := []struct {
		name          string
		args          []string
		expectedError string
	}{
		{
			name:          "Less than 1 row",
			args:          []string{"cmd", "-rows", "0"},
			expectedError: "Invalid flags: invalid number of rows: 0",
		},
		{
			name:          "No fields",
			args:          []string{"cmd", "-fields", ""},
			expectedError: "Invalid flags: fields cannot be empty",
		},
		{
			name:          "No output file name",
			args:          []string{"cmd", "-filename", ""},
			expectedError: "Invalid flags: filename cannot be empty",
		},
		{
			name:          "Invalid fields",
			args:          []string{"cmd", "-fields", "invalid"},
			expectedError: "Unable to generate CSV data. Invalid fields selected: invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			os.Args = tt.args

			defer func() {
				if err := recover(); err != nil {
					if err != tt.expectedError {
						t.Errorf("Expected error: %v\nGot: %v", tt.expectedError, err)
					}
				}
			}()

			main()
		})
	}
}

func TestMain_SuccessCases(t *testing.T) {
	origStdout := os.Stdout
	origArgs := os.Args
	defer func() {
		os.Stdout = origStdout
		os.Args = origArgs

		// os.RemoveAll("output")
	}()

	tests := []struct {
		name             string
		args             []string
		expectedOut      string
		filename         string
		expectedFileData [][]string
	}{
		{
			name:             "Default values",
			args:             []string{"cmd", "-seed", "1"},
			expectedOut:      "CSV file successfully generated at output/output.csv.",
			filename:         "output.csv",
			expectedFileData: [][]string{{"name", "age"}, {"Zion Brakus", "94"}},
		},
		{
			name:             "Two rows",
			args:             []string{"cmd", "-rows", "2", "-seed", "1"},
			expectedOut:      "CSV file successfully generated at output/output.csv.",
			filename:         "output.csv",
			expectedFileData: [][]string{{"name", "age"}, {"Zion Brakus", "94"}, {"Randy Braun", "98"}},
		},
		{
			name:             "Custom fields",
			args:             []string{"cmd", "-fields", "email,firstName,lastName,city", "-seed", "1"},
			expectedOut:      "CSV file successfully generated at output/output.csv.",
			filename:         "output.csv",
			expectedFileData: [][]string{{"email", "firstName", "lastName", "city"}, {"zion.brakus@productparadigms.biz", "Zion", "Brakus", "Irving"}},
		},
		{
			name:             "Custom file name",
			args:             []string{"cmd", "-filename", "test_data.csv", "-seed", "1"},
			expectedOut:      "CSV file successfully generated at output/test_data.csv.",
			filename:         "test_data.csv",
			expectedFileData: [][]string{{"name", "age"}, {"Zion Brakus", "94"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			os.Args = tt.args

			r, w, _ := os.Pipe()
			os.Stdout = w

			main()

			w.Close()
			var buf bytes.Buffer
			io.Copy(&buf, r)

			actualOut := buf.String()
			lines := strings.Split(actualOut, "\n")
			output := lines[4]

			if output != tt.expectedOut {
				t.Errorf("\nExpected output:\n%s\nGot:\n%s", tt.expectedOut, output)
			}

			// Read and verify the output file
			outputFile, err := os.Open(fmt.Sprintf("output/%s", tt.filename))
			if err != nil {
				t.Errorf("Failed to open output file: %v", err)
			}
			defer outputFile.Close()

			reader := csv.NewReader(outputFile)
			records, err := reader.ReadAll()
			if err != nil {
				t.Errorf("Failed to read CSV file: %v", err)
			}

			for idx, record := range records {
				expectedRow := strings.Join(tt.expectedFileData[idx], ",")
				actualRow := strings.Join(record, ",")
				if actualRow != expectedRow {
					t.Errorf("\nFile data mistmatch at index %d.\nExpected:\n%s\nGot:\n%s", idx, expectedRow, actualRow)
				}
			}
		})
	}
}

func TestGenerate_ErrorCases(t *testing.T) {
	rows := 1
	fields := "email"
	filename := "output.csv"
	seed := 1

	tests := []struct {
		name          string
		args          []string
		fileHandler   FileHandler
		fileWriter    FileWriter
		dataGenerator DataGenerator
		expectedError string
	}{
		{
			name:          "Generate csv data fails",
			fileHandler:   &MockFileHandler{},
			fileWriter:    &MockFileWriter{},
			dataGenerator: &MockDataGenerator{ShouldFail: true},
			expectedError: "Failed to generate CSV data: generateCsvData failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					if err != tt.expectedError {
						t.Errorf("Expected error: %v\nGot: %v", tt.expectedError, err)
					}
				}
			}()

			generate(tt.fileHandler, tt.fileWriter, tt.dataGenerator, &rows, &fields, &filename, &seed)
		})
	}
}

func TestGenerate_SuccessCases(t *testing.T) {
	origStdout := os.Stdout
	defer func() {
		os.Stdout = origStdout
	}()

	rows := 1
	fields := "email"
	filename := "output.csv"
	seed := 1

	tests := []struct {
		name          string
		args          []string
		fileHandler   FileHandler
		fileWriter    FileWriter
		dataGenerator DataGenerator
		expectedOut   string
	}{
		{
			name:          "Generate csv data success",
			args:          []string{"cmd", "-rows", "1"},
			fileHandler:   &MockFileHandler{},
			fileWriter:    &MockFileWriter{},
			dataGenerator: &MockDataGenerator{ShouldFail: false},
			expectedOut:   "CSV file successfully generated at output/output.csv.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, w, _ := os.Pipe()
			os.Stdout = w

			generate(tt.fileHandler, tt.fileWriter, tt.dataGenerator, &rows, &fields, &filename, &seed)

			w.Close()
			var buf bytes.Buffer
			io.Copy(&buf, r)

			actualOut := buf.String()
			lines := strings.Split(actualOut, "\n")
			output := lines[4]

			if output != tt.expectedOut {
				t.Errorf("\nExpected output:\n%s\nGot:\n%s", tt.expectedOut, output)
			}
		})
	}
}

func TestGenerateCsvData_ErrorCases(t *testing.T) {
	rows := 1
	fields := "email"
	outputDir := "output"
	filename := "output.csv"
	dataGenerator := CSVDataGenerator{}

	tests := []struct {
		name          string
		args          []string
		fileHandler   FileHandler
		fileWriter    FileWriter
		expectedError string
	}{
		{
			name:          "FileHandler.MkDirAll fails",
			fileHandler:   &MockFileHandler{ShouldFailMkDirAll: true},
			fileWriter:    &MockFileWriter{ShouldFail: false},
			expectedError: "failed to create directory: MkDirAll failed",
		},
		{
			name:          "FilerHandler.Create fails",
			fileHandler:   &MockFileHandler{ShouldFailCreate: true},
			fileWriter:    &MockFileWriter{ShouldFail: false},
			expectedError: "Create failed",
		},
		{
			name:          "FileWriter.Write fails",
			fileHandler:   &MockFileHandler{},
			fileWriter:    &MockFileWriter{ShouldFail: true},
			expectedError: "failed to write header row: Write failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dataGenerator.generateCsvData(rows, fields, outputDir, filename, tt.fileHandler, tt.fileWriter)

			if err == nil || err.Error() != tt.expectedError {
				t.Errorf("Expected error: %v\nGot: %v", tt.expectedError, err)
			}
		})
	}
}

func TestGenerateCsvData_SuccessCases(t *testing.T) {
	rows := 1
	fields := "email"
	outputDir := "output"
	filename := "output.csv"
	dataGenerator := CSVDataGenerator{}

	tests := []struct {
		name        string
		args        []string
		fileHandler FileHandler
		fileWriter  FileWriter
	}{
		{
			name:        "Successfully write to csv data file",
			fileHandler: &MockFileHandler{},
			fileWriter:  &MockFileWriter{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dataGenerator.generateCsvData(rows, fields, outputDir, filename, tt.fileHandler, tt.fileWriter)

			if err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
		})
	}
}
