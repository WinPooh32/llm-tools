package main

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"slices"
	"strconv"
	"strings"

	"github.com/WinPooh32/llm-tools/pkg/mcputil"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type RunCommandInput struct {
	Command string `json:"command" jsonschema:"command with argumets"`
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

	outputStr := string(output)

	if exitCode != 0 {
		outputStr += "\n\nExit Code status: " + strconv.Itoa(exitCode)
	}

	return mcputil.TextResult(outputStr), nil, nil
}
