// RPC客户端逻辑测试
package solana

import (
	"testing"
)

// 测试Helius URL拼接
func TestBuildHeliusURL(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		apiKey   string
		want     string
		wantErr  bool
	}{
		{
			name:     "正常拼接",
			endpoint: "https://mainnet.helius-rpc.com/",
			apiKey:   "test-key-123",
			want:     "https://mainnet.helius-rpc.com/?api-key=test-key-123",
		},
		{
			name:     "endpoint已含?",
			endpoint: "https://mainnet.helius-rpc.com/?other=1",
			apiKey:   "test-key-123",
			want:     "https://mainnet.helius-rpc.com/?other=1&api-key=test-key-123",
		},
		{
			name:     "endpoint已含api-key则原样返回",
			endpoint: "https://mainnet.helius-rpc.com/?api-key=existing",
			apiKey:   "ignored",
			want:     "https://mainnet.helius-rpc.com/?api-key=existing",
		},
		{
			name:     "空endpoint使用默认",
			endpoint: "",
			apiKey:   "test-key-123",
			want:     "https://mainnet.helius-rpc.com/?api-key=test-key-123",
		},
		{
			name:     "无apiKey报错",
			endpoint: "https://mainnet.helius-rpc.com/",
			apiKey:   "",
			wantErr:  true,
		},
		{
			name:     "空白apiKey报错",
			endpoint: "https://mainnet.helius-rpc.com/",
			apiKey:   "   ",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildHeliusURL(tt.endpoint, tt.apiKey)
			if tt.wantErr {
				if err == nil {
					t.Error("期望错误但未返回")
				}
				return
			}
			if err != nil {
				t.Fatalf("不期望错误: %v", err)
			}
			if got != tt.want {
				t.Errorf("buildHeliusURL() = %q, want %q", got, tt.want)
			}
		})
	}
}
