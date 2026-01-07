package main

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/WinPooh32/llm-tools/pkg/mcputil"
	"github.com/WinPooh32/llm-tools/tools/file/lines"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ReadInput struct {
	Path         string `json:"path"                    jsonschema:"path to the file"`
	DisableLines *bool  `json:"disable_lines,omitempty" jsonschema:"disable line numbers"`
}

func Read(ctx context.Context, _ *mcp.CallToolRequest, input ReadInput) (*mcp.CallToolResult, any, error) {
	file, err := openFile(cwd, input.Path, true)
	if errors.Is(err, errEscapesFromParent) {
		return mcputil.ErrorResult(escapesFromParentErr), nil, nil
	}
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	bs, err := io.ReadAll(file)
	if err != nil {
		return nil, nil, fmt.Errorf("read file: io: %w", err)
	}

	if input.DisableLines != nil && *input.DisableLines {
		return mcputil.TextResult(string(bs)), nil, nil
	}

	return mcputil.TextResult(lines.AddNumbers(string(bs))), nil, nil
}
