// Helius索引RPC客户端
package solana

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

// RPC请求结构体
type rpcRequest struct {
	Jsonrpc string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// RPC客户端
type RPCClient struct {
	URL    string
	Client *http.Client
}

// NewRPCClient 创建RPC客户端(Helius)
func NewRPCClient(endpoint, apiKey string, timeout time.Duration) (*RPCClient, error) {
	// 组装URL
	url, err := buildHeliusURL(endpoint, apiKey)
	if err != nil {
		return nil, err
	}
	// 初始化HTTP客户端
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	client := &http.Client{Timeout: timeout}
	return &RPCClient{URL: url, Client: client}, nil
}

// buildHeliusURL 拼接Helius地址
func buildHeliusURL(endpoint, apiKey string) (string, error) {
	// endpoint为空就使用默认
	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		endpoint = "https://mainnet.helius-rpc.com/"
	}

	// endpoint已包含api-key则直接使用
	if strings.Contains(endpoint, "api-key=") {
		return endpoint, nil
	}

	// 需要apiKey
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return "", fmt.Errorf("必须提供 Helius API Key")
	}
	sep := "?"
	if strings.Contains(endpoint, "?") {
		sep = "&"
	}
	return endpoint + sep + "api-key=" + apiKey, nil
}

// Call 发送RPC请求并返回gjson结果
func (c *RPCClient) Call(ctx context.Context, method string, params interface{}) (gjson.Result, error) {
	// 组装请求体
	payload := rpcRequest{Jsonrpc: "2.0", ID: 1, Method: method, Params: params}
	body, err := json.Marshal(payload)
	if err != nil {
		return gjson.Result{}, err
	}

	// 发起请求
	req, err := http.NewRequestWithContext(ctx, "POST", c.URL, bytes.NewReader(body))
	if err != nil {
		return gjson.Result{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "meme-chart/1.0")

	resp, err := c.Client.Do(req)
	if err != nil {
		return gjson.Result{}, err
	}
	defer resp.Body.Close()

	// 处理非200
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return gjson.Result{}, fmt.Errorf("HTTP状态码异常: %d, %s", resp.StatusCode, string(bodyBytes))
	}

	// 读取并解析JSON
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return gjson.Result{}, err
	}
	return gjson.ParseBytes(bodyBytes), nil
}
