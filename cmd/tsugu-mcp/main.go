// Command tsugu-mcpは相続書類PDFを生成するMCPサーバー(stdio)
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/chan-mai/tsugu-mcp/internal/buildinfo"
	"github.com/chan-mai/tsugu-mcp/internal/mcpserver"
)

func main() {
	for _, a := range os.Args[1:] {
		if a == "-version" || a == "--version" || a == "-v" {
			fmt.Println("tsugu-mcp", buildinfo.String())
			return
		}
	}
	if err := mcpserver.Run(context.Background()); err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
