package main

import (
	"testing"
)

func TestGetAnalyzers(t *testing.T) {
	analyzers := getAnalyzers()
	if len(analyzers) == 0 {
		t.Fatal("expected at least one analyzer")
	}
	var found bool
	for _, a := range analyzers {
		if a == NoOsExitAnalyzer {
			found = true
			break
		}
	}
	if !found {
		t.Error("NoOsExitAnalyzer not found in returned analyzers")
	}
}
