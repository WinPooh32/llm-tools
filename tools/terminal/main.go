package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type argList []string

func (a *argList) String() string {
	return strings.Join(*a, ", ")
}

func (a *argList) Set(s string) error {
	*a = strings.Split(s, ",")

	for i, v := range *a {
		(*a)[i] = strings.TrimSpace(v)
	}

	return nil
}

var allowedCommands argList

func main() {
	httpAddr := flag.String("http", "", "if set, use streamable HTTP at this address, instead of stdin/stdout")
	flag.Var(&allowedCommands, "allowed-commands", "comma-separated list of allowed commands")
	flag.Parse()

	server := mcp.NewServer(&mcp.Implementation{Name: "file"}, nil)

	var description string

	if len(allowedCommands) > 0 {
		description = "Exec one of these allowed commands: " + allowedCommands.String()
	} else {
		description = "Exec command"
	}

	mcp.AddTool(server,
		&mcp.Tool{
			Name:        "terminal",
			Description: description,
		},
		Terminal,
	)

	if *httpAddr != "" {
		handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
			return server
		}, nil)

		log.Printf("File MCP server listening at %s", *httpAddr)
		if err := http.ListenAndServe(*httpAddr, handler); err != nil {
			log.Fatal(err)
		}
	} else {
		t := &mcp.LoggingTransport{Transport: &mcp.StdioTransport{}, Writer: os.Stderr}
		if err := server.Run(context.Background(), t); err != nil {
			log.Printf("Server failed: %v", err)
		}
	}
}
