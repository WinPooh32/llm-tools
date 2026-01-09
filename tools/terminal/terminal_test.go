package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_explodeCommand(t *testing.T) {
	tests := []struct {
		name    string
		command string
		want    []string
	}{
		{"simple", `go run main.go`, []string{"go", "run", "main.go"}},
		{"quoted", `go run "main.go"`, []string{"go", "run", "main.go"}},
		{"quoted with spaces", `bash -c "go run main.go"`, []string{"bash", "-c", "go run main.go"}},
		{"quoted with spaces and escaped quotes", `bash -c "go run \"main.go\""`, []string{"bash", "-c", `go run "main.go"`}},
		{"quited bash ls with subquotes", `bash -c "ls \".\""`, []string{"bash", "-c", `ls "."`}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := explodeCommand(tt.command)
			assert.Equal(t, tt.want, got)
		})
	}
}
