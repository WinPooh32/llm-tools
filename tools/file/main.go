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

Line parameters are 1-indexed.

Examples:
- To insert a new line at line 5: begin_line=5, end_line=5, content="new line content\n"
- To replace lines 3-7 with new content: begin_line=3, end_line=8, content="new line 1\nnew line 2\n..."
- To delete lines 10-15: begin_line=10, end_line=16, content=null
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
