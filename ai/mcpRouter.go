package ai

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"xhhrobot/config"
	"xhhrobot/loger"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"
)

type ServerSession struct {
	name    string
	session *mcp.ClientSession
}

type MultiMCPServerManager struct {
	mu             sync.RWMutex
	sessions       map[string]*ServerSession
	client         *mcp.Client
	cachedToolDefs []ToolDef
}

var mcpMgr *MultiMCPServerManager

func Init() {
	{
		cfg := config.ConfigStruct.Ai.MCP
		if cfg.Enabled {
			mcpMgr = newMultiMCPServerManager()
			for serverName, v := range cfg.MCPServers {
				cmd := exec.Command(v.Command, v.Args...)
				if cfg.UseOSEnv {
					cmd.Env = os.Environ()
				}
				for key, val := range v.Env {
					cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, val))
				}
				if err := mcpMgr.connect(context.Background(), serverName, &mcp.CommandTransport{Command: cmd}); err != nil {
					loger.Loger.Warn("[AI]fail to connect mcp server", zap.String("mcpServerName", serverName))
				}
			}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			mcpTools, err := mcpMgr.listAllTools(ctx)
			cancel()
			if err != nil {
				loger.Loger.Error("[AI]fail to list mcp tools")
			}
			mcpMgr.cachedToolDefs = mcpToToolDefs(mcpTools)
			loger.Loger.Info("[AI]cached mcp tools", zap.Int("toolsCnt", len(mcpMgr.cachedToolDefs)))
		}
	}

	loger.Loger.Info("[AI]Init")
}

func Close() {
	if mcpMgr != nil {
		mcpMgr.close()
	}
}

func newMultiMCPServerManager() *MultiMCPServerManager {
	return &MultiMCPServerManager{
		sessions: make(map[string]*ServerSession),
		client: mcp.NewClient(&mcp.Implementation{
			Name:    "muti-servers-manager",
			Version: "1.0",
		}, nil),
	}
}

func (m *MultiMCPServerManager) connect(ctx context.Context, serverName string, transport mcp.Transport) error {
	session, err := m.client.Connect(ctx, transport, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", serverName, err)
	}

	m.mu.Lock()
	m.sessions[serverName] = &ServerSession{name: serverName, session: session}
	m.mu.Unlock()

	loger.Loger.Info("[AI]connected to server", zap.String("serverName", serverName))
	return nil
}

// 聚合所有 Server 的工具，并添加前缀防止重名
func (m *MultiMCPServerManager) listAllTools(ctx context.Context) ([]mcp.Tool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var allTools []mcp.Tool
	for name, sess := range m.sessions {
		var curCursor string
		for {
			page, err := sess.session.ListTools(ctx, &mcp.ListToolsParams{Cursor: curCursor})
			if err != nil {
				return nil, fmt.Errorf("list tools from %s failed: %w", name, err)
			}

			// 双下划线前缀
			for _, tool := range page.Tools {
				newTool := *tool
				newTool.Name = fmt.Sprintf("%s__%s", name, tool.Name)
				allTools = append(allTools, newTool)
			}

			if page.NextCursor == "" {
				break
			}
			curCursor = page.NextCursor
		}
	}
	return allTools, nil
}

// 解析前缀，将请求路由到正确的 Server
func (m *MultiMCPServerManager) callTool(ctx context.Context, prefixedToolName string, args map[string]any) (*mcp.CallToolResult, error) {
	parts := strings.SplitN(prefixedToolName, "__", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid tool name format, expected 'server__tool'")
	}
	serverName, realToolName := parts[0], parts[1]

	m.mu.RLock()
	sess, ok := m.sessions[serverName]
	m.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("server %s not found or not connected", serverName)
	}

	t := config.ConfigStruct.Ai.MCP.ToolCallTimeLimit
	if t <= 0 {
		t = 30
	}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(t)*time.Second)
	res, err := sess.session.CallTool(ctx, &mcp.CallToolParams{Name: realToolName, Arguments: args})
	cancel()

	return res, err
}

func (m *MultiMCPServerManager) close() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, sess := range m.sessions {
		sess.session.Close()
	}
}

func mcpToToolDefs(mcpTools []mcp.Tool) []ToolDef {
	defs := make([]ToolDef, 0, len(mcpTools))
	for _, t := range mcpTools {
		params, _ := t.InputSchema.(map[string]any)
		defs = append(defs, ToolDef{
			Type: "function",
			Function: ToolFunc{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  params,
			},
		})
	}
	return defs
}

func extractToolResult(result *mcp.CallToolResult, err error) string {
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	if result == nil {
		return "Error: nil result"
	}
	var sb strings.Builder
	for _, c := range result.Content {
		if tc, ok := c.(*mcp.TextContent); ok {
			sb.WriteString(tc.Text)
		}
	}
	if result.IsError && sb.Len() == 0 {
		return "Error: tool call failed with no content"
	}
	return sb.String()
}
