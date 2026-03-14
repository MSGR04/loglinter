package analyzer

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Current directory: %s", wd)

	testdata := filepath.Join(wd, "testdata")
	t.Logf("Looking for testdata at: %s", testdata)

	if _, err := os.Stat(testdata); os.IsNotExist(err) {
		t.Fatalf("Testdata directory does not exist: %s", testdata)
	}
	testFile := filepath.Join(testdata, "src", "mylogger", "mylogger.go")
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Fatalf("Test file does not exist: %s", testFile)
	}
	analysistest.Run(t, testdata, Analyzer, "mylogger")
}
