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
	Path    string `json:"path"       jsonschema:"path to the file"`
	Begin   int    `json:"begin_line" jsonschema:"begin line number of the edit selection"`
	End     int    `json:"end_line"   jsonschema:"end line number of the edit selection"`
	Content string `json:"content"    jsonschema:"text content that will be applied to the selected region"`
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
	if input.Content != nil {
		editLines = strings.Split(*input.Content, "\n")
	}

	newLines := applyLines(srcLines, editLines, input.Begin-1, input.End-1)

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

	return mcputil.TextResult("OK"), nil, nil
}

func applyLines(lines, newLines []string, begin, end int) []string {
	if end < begin {
		return lines
	}

	if begin >= len(lines) {
		begin = len(lines)
	}

	if end > len(lines) {
		end = len(lines)
	}

	return slices.Concat(lines[:begin], newLines, lines[end:])
}
