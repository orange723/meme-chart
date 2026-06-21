// 静态HTML渲染
package render

import (
	"encoding/json"
	"html/template"
	"os"

	"meme-chart/internal/model"
)

// RenderStaticHTML 渲染静态HTML
func RenderStaticHTML(holders []model.Holder, meta model.TokenMeta, chartMode string, outPath string) error {
	// 组装JSON数据
	payload := map[string]interface{}{
		"meta":      meta,
		"holders":   holders,
		"chartMode": chartMode,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// 解析模板
	tmpl := template.Must(template.New("index").Parse(DynamicHTML))

	// 构建页面数据（Interval=0 表示静态模式）
	data := PageData{
		Title:       BuildTitle(meta),
		Meta:        meta,
		ChartMode:   chartMode,
		PayloadJSON: template.JS(payloadBytes),
	}

	// 输出HTML
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}
