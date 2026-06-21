// 渲染逻辑测试
package render

import (
	"encoding/json"
	"html/template"
	"os"
	"strings"
	"testing"

	"meme-chart/internal/model"
)

// 测试BuildTitle
func TestBuildTitle(t *testing.T) {
	tests := []struct {
		name     string
		meta     model.TokenMeta
		expected string
	}{
		{
			name:     "名称+符号",
			meta:     model.TokenMeta{Name: "Dogwifhat", Symbol: "WIF", Mint: "EKpQGSJtjMFqKZ9KQanSqYXRcF8fBopzLHYxdM65zcjm"},
			expected: "Dogwifhat (WIF)",
		},
		{
			name:     "仅有名称",
			meta:     model.TokenMeta{Name: "Dogwifhat", Symbol: "", Mint: "EKpQGSJtjMFqKZ9KQanSqYXRcF8fBopzLHYxdM65zcjm"},
			expected: "Dogwifhat",
		},
		{
			name:     "仅有Symbol",
			meta:     model.TokenMeta{Name: "", Symbol: "WIF", Mint: "EKpQGSJtjMFqKZ9KQanSqYXRcF8fBopzLHYxdM65zcjm"},
			expected: "EKpQGSJtjMFqKZ9KQanSqYXRcF8fBopzLHYxdM65zcjm",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildTitle(tt.meta)
			if got != tt.expected {
				t.Errorf("BuildTitle() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// 测试静态HTML渲染
func TestRenderStaticHTML(t *testing.T) {
	meta := model.TokenMeta{
		Name:     "Test Token",
		Symbol:   "TEST",
		Mint:     "test_mint_address",
		Supply:   1_000_000.0,
		Decimals: 6,
	}

	holders := []model.Holder{
		{Owner: "alice...1111", UiAmount: 500_000.0},
		{Owner: "bob...2222", UiAmount: 300_000.0},
		{Owner: "carol...3333", UiAmount: 200_000.0},
	}

	// 写入临时文件
	tmpFile := t.TempDir() + "/test_output.html"

	err := RenderStaticHTML(holders, meta, "pie", tmpFile)
	if err != nil {
		t.Fatalf("RenderStaticHTML 失败: %v", err)
	}

	// 读取并验证
	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("读取输出文件失败: %v", err)
	}
	html := string(content)

	// 验证关键内容存在
	checks := []string{
		"<!DOCTYPE html>",
		"Test Token (TEST)",
		"echarts",
		"alice...1111",
		"bob...2222",
		"carol...3333",
		"seriesData",
	}
	for _, check := range checks {
		if !strings.Contains(html, check) {
			t.Errorf("输出HTML中缺少: %q", check)
		}
	}
}

// 测试模板解析（确保 DynamicHTML 可解析）
func TestTemplateParsing(t *testing.T) {
	tmpl, err := template.New("test").Parse(DynamicHTML)
	if err != nil {
		t.Fatalf("模板解析失败: %v", err)
	}

	// 测试静态模式渲染
	meta := model.TokenMeta{Name: "Test", Symbol: "TST", Mint: "mint123", Supply: 100.0}
	data := PageData{
		Title:       "Test (TST)",
		Meta:        meta,
		ChartMode:   "pie",
		Interval:    0,
		PayloadJSON: template.JS(`{"holders":[],"chartMode":"pie"}`),
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		t.Fatalf("模板执行失败: %v", err)
	}
	if !strings.Contains(buf.String(), "StaticPayloadJSON") && !strings.Contains(buf.String(), "renderChart") {
		t.Error("模板产出缺少图表渲染代码")
	}

	// 验证JSON内嵌（静态模式不包含 loadData fetch 逻辑）
	if strings.Contains(buf.String(), "fetch('/data')") {
		t.Error("静态模式下不应包含 fetch('/data') 代码")
	}

	// 测试动态模式渲染
	data2 := PageData{
		Title:     "Test (TST)",
		Meta:      meta,
		ChartMode: "pie",
		Interval:  30,
	}
	var buf2 strings.Builder
	if err := tmpl.Execute(&buf2, data2); err != nil {
		t.Fatalf("模板执行失败: %v", err)
	}
	// 动态模式应包含刷新文字和 fetch 调用
	if !strings.Contains(buf2.String(), "刷新: 30 秒") {
		t.Error("动态模式下应显示刷新间隔")
	}
	if !strings.Contains(buf2.String(), "fetch('/data')") {
		t.Error("动态模式下应包含 fetch('/data') 调用")
	}
}

// 测试模板 PayloadJSON 正确序列化
func TestPayloadJSONSerialization(t *testing.T) {
	payload := map[string]interface{}{
		"holders": []map[string]interface{}{
			{"owner": "test...1111", "uiAmount": 100.0},
		},
		"chartMode": "pie",
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}
	payloadJS := template.JS(payloadBytes)

	meta := model.TokenMeta{Name: "Test", Symbol: "TST", Mint: "mint123", Supply: 100.0}
	data := PageData{
		Title:       "Test (TST)",
		Meta:        meta,
		ChartMode:   "pie",
		PayloadJSON: payloadJS,
	}

	tmpl := template.Must(template.New("test").Parse(DynamicHTML))
	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		t.Fatalf("模板执行失败: %v", err)
	}
	if !strings.Contains(buf.String(), `"test...1111"`) {
		t.Error("PayloadJSON未正确嵌入HTML")
	}
}
