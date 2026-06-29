// Command tsugu-mcpは相続書類PDFを生成するMCPサーバー(stdio)
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/chan-mai/tsugu-mcp/internal/mcpserver"
)

func main() {
	if err := mcpserver.Run(context.Background()); err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
