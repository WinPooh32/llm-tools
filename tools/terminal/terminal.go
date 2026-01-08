package main

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"slices"
	"strings"

	"github.com/WinPooh32/llm-tools/pkg/mcputil"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type RunCommandInput struct {
	Command string `json:"command" jsonschema:"command with argumets"`
}

type RunCommandOutput struct {
	ExitStatus int    `json:"exit_status" jsonschema:"exit status code of the command"`
	Output     string `json:"output"      jsonschema:"command output"`
}

func RunCommand(ctx context.Context, _ *mcp.CallToolRequest, input RunCommandInput) (*mcp.CallToolResult, any, error) {
	command := strings.TrimSpace(input.Command)

	if len(command) == 0 {
		return mcputil.ErrorResult("failed to exec empty command"), nil, nil
	}

	parts := strings.Split(command, " ")

	for i, v := range parts {
		parts[i] = strings.TrimSpace(v)
	}

	var (
		argc string
		argv []string
	)

	if len(parts) == 0 {
		return mcputil.ErrorResult("failed to exec empty command"), nil, nil
	}

	argc = parts[0]

	if len(parts) > 1 {
		argv = parts[1:]
	}

	if len(allowedCommands) > 0 && !slices.Contains(allowedCommands, argc) {
		return mcputil.ErrorResult(`command "` + argc + `" is not allowed`), nil, nil
	}

	cmd := exec.CommandContext(ctx, argc, argv...)
	output, err := cmd.CombinedOutput()

	var exitErr *exec.ExitError
	if err != nil && !errors.As(err, &exitErr) {
		return nil, nil, fmt.Errorf("failed to exec command: %w", err)
	}

	exitCode := 0
	if cmd.ProcessState != nil {
		exitCode = cmd.ProcessState.ExitCode()
	}

	return nil, &RunCommandOutput{
		ExitStatus: exitCode,
		Output:     string(output),
	}, nil
}
