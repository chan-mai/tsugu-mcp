package mcpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/chan-mai/tsugu-mcp/registration"
	"github.com/chan-mai/tsugu-mcp/relationchart"
)

// era文字列を表記スタイルへ変換(chart/touki共通)
func parseEra(s string) (relationchart.EraStyle, error) {
	switch s {
	case "", "wareki":
		return relationchart.EraWareki, nil
	case "both":
		return relationchart.EraWarekiWithSeireki, nil
	case "seireki":
		return relationchart.EraSeireki, nil
	default:
		return relationchart.EraWareki, fmt.Errorf("unknown era %q (wareki|both|seireki)", s)
	}
}

// PDFをファイルへ書き出し絶対パスを返却(outputPath空は一時ファイル)
func writePDF(pdf []byte, outputPath, kind string) (string, error) {
	if outputPath == "" {
		f, err := os.CreateTemp("", "tsugu-"+kind+"-*.pdf")
		if err != nil {
			return "", err
		}
		defer f.Close()
		if _, err := f.Write(pdf); err != nil {
			return "", err
		}
		return f.Name(), nil
	}
	if dir := filepath.Dir(outputPath); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return "", err
		}
	}
	if err := os.WriteFile(outputPath, pdf, 0o644); err != nil {
		return "", err
	}
	if abs, err := filepath.Abs(outputPath); err == nil {
		return abs, nil
	}
	return outputPath, nil
}

func errorResult(err error) (*mcp.CallToolResult, toolResult, error) {
	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{&mcp.TextContent{Text: err.Error()}},
	}, toolResult{}, nil
}

func okResult(path string, n int) (*mcp.CallToolResult, toolResult, error) {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("generated: %s (%d bytes)", path, n)}},
	}, toolResult{Path: path, Bytes: n}, nil
}

func handleChart(_ context.Context, _ *mcp.CallToolRequest, in chartToolInput) (*mcp.CallToolResult, toolResult, error) {
	opt, err := parseEra(in.Era)
	if err != nil {
		return errorResult(err)
	}
	data, err := json.Marshal(in.Document)
	if err != nil {
		return errorResult(err)
	}
	pdf, err := relationchart.GenerateFromJSON(data, relationchart.Options{Era: opt})
	if err != nil {
		return errorResult(err)
	}
	path, err := writePDF(pdf, in.OutputPath, "chart")
	if err != nil {
		return errorResult(err)
	}
	return okResult(path, len(pdf))
}

func handleTouki(_ context.Context, _ *mcp.CallToolRequest, in toukiToolInput) (*mcp.CallToolResult, toolResult, error) {
	opt, err := parseEra(in.Era)
	if err != nil {
		return errorResult(err)
	}
	taxNote := autoFillTax(&in.Document)
	data, err := json.Marshal(in.Document)
	if err != nil {
		return errorResult(err)
	}
	pdf, err := registration.GenerateFromJSON(data, registration.Options{Era: opt})
	if err != nil {
		return errorResult(err)
	}
	path, err := writePDF(pdf, in.OutputPath, "touki")
	if err != nil {
		return errorResult(err)
	}
	res, out, e := okResult(path, len(pdf))
	if taxNote != "" {
		res.Content = append(res.Content, &mcp.TextContent{Text: taxNote})
	}
	return res, out, e
}
