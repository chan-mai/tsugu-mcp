package mcpserver

import "github.com/chan-mai/tsugu-mcp/docs/knowledge"

// kbOutputは知識系ツールの構造化出力に調査基準日を添えるラッパ
type kbOutput[T any] struct {
	AsOf   string `json:"asOf" jsonschema:"知識ベースの調査基準日 YYYY-MM-DD。この日時点の法令・税率を反映"`
	Result T      `json:"result" jsonschema:"算定・案内の結果"`
}

func withAsOf[T any](r T) kbOutput[T] { return kbOutput[T]{AsOf: knowledge.AsOf, Result: r} }
