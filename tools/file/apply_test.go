package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_applyLines(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		lines    []string
		newLines []string
		begin    int
		end      int
		want     []string
	}{
		{
			name:     "insert between lines",
			lines:    []string{"1", "2", "3"},
			newLines: []string{"2.5"},
			begin:    2,
			end:      2,
			want:     []string{"1", "2", "2.5", "3"},
		},
		{
			name:     "insert at end",
			lines:    []string{"1", "2", "3"},
			newLines: []string{"4"},
			begin:    3,
			end:      3,
			want:     []string{"1", "2", "3", "4"},
		},
		{
			name:     "insert at begin",
			lines:    []string{"1", "2", "3"},
			newLines: []string{"0"},
			begin:    0,
			end:      0,
			want:     []string{"0", "1", "2", "3"},
		},
		{
			name:     "replace multiple lines at end",
			lines:    []string{"1", "2", "3", "4"},
			newLines: []string{"2.1", "2.2"},
			begin:    2,
			end:      4,
			want:     []string{"1", "2", "2.1", "2.2"},
		},
		{
			name:     "replace multiple lines at begin",
			lines:    []string{"1", "2", "3", "4"},
			newLines: []string{"0.1", "0.2"},
			begin:    0,
			end:      2,
			want:     []string{"0.1", "0.2", "3", "4"},
		},
		{
			name:     "remove line at the middle",
			lines:    []string{"1", "2", "3"},
			newLines: []string{},
			begin:    1,
			end:      2,
			want:     []string{"1", "3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := applyLines(tt.lines, tt.newLines, tt.begin, tt.end)
			assert.Equal(t, tt.want, got)
		})
	}
}
