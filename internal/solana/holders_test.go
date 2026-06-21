// 持仓数据处理逻辑测试
package solana

import (
	"testing"

	"meme-chart/internal/model"
)

// 测试地址缩短
func TestShortenAddress(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"短地址原样返回", "abc", "abc"},
		{"10字符边界", "1234567890", "1234567890"},
		{"标准Solana地址", "7xKXtg2CW87dQL4dJdPxPCz9eqh7oF8nNdHh5v4pump", "7xKX...pump"},
		{"刚好12字符", "123456789012", "1234...9012"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shortenAddress(tt.input)
			if got != tt.expected {
				t.Errorf("shortenAddress(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// 测试排序与截断
func TestSortHolders(t *testing.T) {
	// 准备数据
	ownerAmounts := map[string]float64{
		"alice1111111111111111111111111111111111111111":  50.0,
		"bob111111111111111111111111111111111111111111":  30.0,
		"carol11111111111111111111111111111111111111111": 20.0,
		"dave111111111111111111111111111111111111111111": 10.0,
	}

	t.Run("Top3", func(t *testing.T) {
		got := sortHolders(ownerAmounts, 3)
		if len(got) != 3 {
			t.Fatalf("len = %d, want 3", len(got))
		}
		// 验证降序
		if got[0].UiAmount < got[1].UiAmount {
			t.Error("结果未按降序排列")
		}
		if got[1].UiAmount < got[2].UiAmount {
			t.Error("结果未按降序排列")
		}
	})

	t.Run("Top超过总数时取全部", func(t *testing.T) {
		got := sortHolders(ownerAmounts, 100)
		if len(got) != 4 {
			t.Fatalf("len = %d, want 4", len(got))
		}
	})

	t.Run("Top为0时不截断", func(t *testing.T) {
		got := sortHolders(ownerAmounts, 0)
		if len(got) != 4 {
			t.Fatalf("len = %d, want 4", len(got))
		}
	})
}

// 测试追加Others
func TestWithOthersBySupply(t *testing.T) {
	holders := []model.Holder{
		{Owner: "alice...1111", UiAmount: 50.0},
		{Owner: "bob...1111", UiAmount: 30.0},
	}

	t.Run("使用总供应量计算Others", func(t *testing.T) {
		got := WithOthersBySupply(holders, 100.0, 80.0)
		if len(got) != 3 {
			t.Fatalf("len = %d, want 3", len(got))
		}
		others := got[len(got)-1]
		if others.Owner != "Others" {
			t.Errorf("Owner = %q, want Others", others.Owner)
		}
		if others.UiAmount != 20.0 {
			t.Errorf("UiAmount = %f, want 20.0", others.UiAmount)
		}
	})

	t.Run("supply为0时回退使用totalFromAccounts", func(t *testing.T) {
		got := WithOthersBySupply(holders, 0, 100.0)
		others := got[len(got)-1]
		if others.UiAmount != 20.0 {
			t.Errorf("UiAmount = %f, want 20.0 (100 - 80)", others.UiAmount)
		}
	})

	t.Run("Others为负数时归零", func(t *testing.T) {
		got := WithOthersBySupply(holders, 60.0, 60.0)
		others := got[len(got)-1]
		if others.UiAmount != 0.0 {
			t.Errorf("UiAmount = %f, want 0.0", others.UiAmount)
		}
	})
}

// 测试 parseFloat 兜底解析
func TestParseFloat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
		hasErr   bool
	}{
		{"纯整数", "12345", 12345, false},
		{"带非法字符", "123abc", 0, true},
		{"空字符串", "", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFloat(tt.input)
			if tt.hasErr && err == nil {
				t.Error("期望错误但未返回")
			}
			if !tt.hasErr && err != nil {
				t.Errorf("不期望错误但返回了: %v", err)
			}
			if got != tt.expected {
				t.Errorf("parseFloat(%q) = %f, want %f", tt.input, got, tt.expected)
			}
		})
	}
}
