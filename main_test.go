package main

import (
	"bytes"
	"flag"
	"io"
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	origStdout := os.Stdout
	origArgs := os.Args
	defer func() {
		os.Stdout = origStdout
		os.Args = origArgs
	}()

	// Test cases
	tests := []struct {
		name        string
		args        []string
		expectedOut string
	}{
		{
			name:        "Default values",
			args:        []string{"cmd"},
			expectedOut: "Rows: 1\nFields: name,age\nFilename: output.csv\n",
		},
		{
			name:        "One-hundred rows",
			args:        []string{"cmd", "-rows", "100"},
			expectedOut: "Rows: 100\nFields: name,age\nFilename: output.csv\n",
		},
		{
			name:        "Custom fields",
			args:        []string{"cmd", "-fields", "email,location"},
			expectedOut: "Rows: 1\nFields: email,location\nFilename: output.csv\n",
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
			if actualOut != tt.expectedOut {
				t.Errorf("Expected output:\n%s\nGot:\n%s", tt.expectedOut, actualOut)
			}
		})
	}
}
