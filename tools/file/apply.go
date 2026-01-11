package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/WinPooh32/llm-tools/pkg/mcputil"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ApplyInput struct {
	Path      string        `json:"path"     jsonschema:"path to the file"`
	Postition ApplyPosition `json:"position" jsonschema:"edit position"`
	Content   *string       `json:"content"  jsonschema:"text content"`
}

func (ai *ApplyInput) IsDelete() bool {
	return ai.Content == nil
}

func (ai *ApplyInput) IsInsert() bool {
	return ai.Postition.End == nil
}

func (ai *ApplyInput) IsReplace() bool {
	return ai.Postition.End != nil
}

func (ai *ApplyInput) Begin() int {
	return ai.Postition.Begin - 1
}

func (ai *ApplyInput) End() int {
	if ai.Postition.End == nil {
		return -1
	}
	return *ai.Postition.End
}

type ApplyPosition struct {
	Begin int  `json:"begin"         jsonschema:"begin line number of the range"`
	End   *int `json:"end,omitempty" jsonschema:"end (inclusive) line number of the range"`
}

func Apply(ctx context.Context, _ *mcp.CallToolRequest, input ApplyInput) (*mcp.CallToolResult, any, error) {
	file, err := openFile(cwd, input.Path, false)
	if errors.Is(err, errEscapesFromParent) {
		return mcputil.ErrorResult(escapesFromParentErr), nil, nil
	}
	if err != nil {
		return nil, nil, err
	}

	bs, err := io.ReadAll(file)
	if err != nil {
		return nil, nil, fmt.Errorf("read file: io: %w", err)
	}

	srcLines := strings.Split(string(bs), "\n")

	var editLines []string
	if input.Content != nil && len(*input.Content) > 0 {
		editLines = strings.Split(*input.Content, "\n")
	}

	var newLines []string

	switch {
	case input.IsDelete():
		newLines = applyLines(srcLines, nil, input.Begin(), input.End())
	case input.IsInsert():
		newLines = applyLines(srcLines, editLines, input.Begin(), -1)
	case input.IsReplace():
		newLines = applyLines(srcLines, editLines, input.Begin(), input.End())
	}

	if err := file.Truncate(0); err != nil {
		return nil, nil, fmt.Errorf("truncate file: %w", err)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, nil, fmt.Errorf("seek to file begin: %w", err)
	}

	if _, err := file.WriteString(strings.Join(newLines, "\n")); err != nil {
		return nil, nil, fmt.Errorf("write to file: %w", err)
	}

	if err := file.Close(); err != nil {
		return nil, nil, fmt.Errorf("close file: %w", err)
	}

	return mcputil.TextResult("SUCCESS\nLine numbers have been changed!"), nil, nil
}

func applyLines(lines, newLines []string, begin, end int) []string {
	if begin < 0 {
		return lines
	}
	if begin > len(lines) {
		begin = len(lines)
	}
	if end > len(lines) || end < 0 {
		end = len(lines)
	}
	if end < begin {
		return lines
	}

	return slices.Concat(lines[:begin], newLines, lines[end:])
}
