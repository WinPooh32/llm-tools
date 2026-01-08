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
If you want to insert content before specific line, begin and end lines must be equal.`,
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
