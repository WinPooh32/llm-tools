package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/WinPooh32/llm-tools/tools/file/lines"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	escapesFromParentErr = "can only access files and directories beneath the current working directory"
)

type ReadInput struct {
	Path         string `json:"path"                    jsonschema:"path to the file"`
	DisableLines *bool  `json:"disable_lines,omitempty" jsonschema:"disable line numbers"`
}

func Read(ctx context.Context, _ *mcp.CallToolRequest, input ReadInput) (*mcp.CallToolResult, any, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, nil, fmt.Errorf("get current working directory: os: %w", err)
	}

	var rel string

	if filepath.IsLocal(input.Path) {
		rel = input.Path
	} else {
		rel, err = filepath.Rel(cwd, input.Path)
		if err != nil {
			return mcpErrorResult(escapesFromParentErr), nil, nil
		}
	}

	file, err := os.OpenInRoot(cwd, rel)
	if os.IsNotExist(err) {
		return mcpErrorResult("file is not exist"), nil, nil
	}
	if err != nil {
		if strings.Contains(err.Error(), "path escapes from parent") {
			return mcpErrorResult(escapesFromParentErr), nil, nil
		}

		return nil, nil, fmt.Errorf("open file: os.Root: %w", err)
	}

	bs, err := io.ReadAll(file)
	if err != nil {
		return nil, nil, fmt.Errorf("read file: io: %w", err)
	}

	if input.DisableLines != nil && *input.DisableLines {
		return mcpTextResult(string(bs)), nil, nil
	}

	return mcpTextResult(lines.AddNumbers(string(bs))), nil, nil
}
