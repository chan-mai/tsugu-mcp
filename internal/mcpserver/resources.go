package mcpserver

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"tsugu-mcp/docs/knowledge"
)

const knowledgeScheme = "knowledge://"

// 知識文書をMCPリソースとして登録
func addKnowledgeResources(s *mcp.Server) {
	for _, d := range knowledge.Docs {
		uri := knowledgeScheme + d.Slug
		s.AddResource(&mcp.Resource{
			URI:         uri,
			Name:        d.Name,
			Description: d.Description,
			MIMEType:    "text/markdown",
		}, knowledgeHandler(uri, d.File))
	}
}

func knowledgeHandler(uri, file string) mcp.ResourceHandler {
	return func(_ context.Context, _ *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		data, err := knowledge.FS.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("read knowledge %q: %w", file, err)
		}
		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{{URI: uri, MIMEType: "text/markdown", Text: string(data)}},
		}, nil
	}
}
