package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	httpAddr := flag.String("http", "", "if set, use streamable HTTP at this address, instead of stdin/stdout")
	flag.Parse()

	server := mcp.NewServer(&mcp.Implementation{Name: "file"}, nil)

	mcp.AddTool(server,
		&mcp.Tool{
			Name: "read",
			Description: `Read a file.
Outputs text where every line has prefixed by it's number.
Format: "<line_number>:<line_content>\n"
Example:
1:first text line
2:second text line
`,
		},
		Read,
	)

	mcp.AddTool(server,
		&mcp.Tool{
			Name: "apply",
			Description: `Apply changes to a text file.

Line numbers are 1-indexed (first line is line 1).

Usage examples:
- Insert a new line at line 5: content: "new line content", position: {begin: 5, end: 5}
- Replace lines 3-7 with new text: content: "new line 1\nnew line 2", position: {begin: 3, end: 7}
- Delete lines 10-15: content: null, position: {begin: 10, end: 15}
`,
		},
		Apply,
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
