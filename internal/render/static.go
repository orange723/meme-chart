// 静态HTML渲染
package render

import (
	"encoding/json"
	"html/template"
	"os"

	"meme-chart/internal/model"
)

// 页面模板数据
type pageData struct {
	Title       string
	Meta        model.TokenMeta
	ChartMode   string
	PayloadJSON template.JS
}

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
	tmpl := template.Must(template.New("index").Parse(dynamicHTML))
	data := pageData{
		Title:       buildTitle(meta),
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

// buildTitle 构建标题
func buildTitle(meta model.TokenMeta) string {
	// 优先使用名称+符号
	if meta.Name != "" && meta.Symbol != "" {
		return meta.Name + " (" + meta.Symbol + ")"
	}
	if meta.Name != "" {
		return meta.Name
	}
	return meta.Mint
}

// 动态页面模板
const dynamicHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>{{.Title}}</title>
  <style>
    body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; margin: 0; background: #f6f7fb; color: #111; }
    header { padding: 18px 24px; background: #111; color: #fff; display: flex; align-items: center; gap: 16px; }
    header img { width: 44px; height: 44px; border-radius: 50%; object-fit: cover; background: #222; }
    .meta { display: flex; flex-direction: column; }
    .meta .title { font-size: 18px; font-weight: 600; }
    .meta .sub { font-size: 12px; color: #bbb; margin-top: 4px; }
    #chart { width: 100%; height: calc(100vh - 80px); }
  </style>
</head>
<body>
  <header>
    {{if .Meta.Image}}<img src="{{.Meta.Image}}" alt="token" />{{end}}
    <div class="meta">
      <div class="title">{{.Title}}</div>
      <div class="sub">总量: {{printf "%.6f" .Meta.Supply}} | 合约: {{.Meta.Mint}}</div>
    </div>
  </header>
  <div id="chart"></div>
  <script src="https://cdn.jsdelivr.net/npm/echarts@5/dist/echarts.min.js"></script>
  <script>
    const payload = {{.PayloadJSON}};
    const holders = payload.holders || [];
    const chartMode = payload.chartMode || '{{.ChartMode}}';
    const seriesData = holders.map(h => ({ name: h.owner, value: h.uiAmount }));
    const chart = echarts.init(document.getElementById('chart'));
    if (chartMode === 'bubble') {
      const max = Math.max(...seriesData.map(s => s.value), 1);
      const palette = ['#1f77b4','#ff7f0e','#2ca02c','#d62728','#9467bd','#8c564b','#e377c2','#7f7f7f','#bcbd22','#17becf'];
      const nodes = seriesData.map((s, i) => ({
        name: s.name,
        value: s.value,
        symbolSize: Math.max(12, Math.sqrt(s.value / max) * 80),
        itemStyle: { color: palette[i % palette.length] }
      }));
      chart.setOption({
        title: { text: 'Holder Bubble Map', left: 'center' },
        tooltip: { formatter: '{b}: {c}' },
        series: [{
          type: 'graph',
          layout: 'force',
          data: nodes,
          roam: true,
          force: { repulsion: 120, edgeLength: 30 },
          label: { show: true, formatter: '{b}' }
        }]
      });
    } else if (chartMode === 'donut') {
      chart.setOption({
        title: { text: 'Holder Distribution', left: 'center' },
        tooltip: { trigger: 'item' },
        legend: { top: 'bottom' },
        series: [{
          type: 'pie',
          radius: ['40%', '70%'],
          center: ['50%', '50%'],
          data: seriesData,
          label: { formatter: '{b}: {d}%'}
        }]
      });
    } else if (chartMode === 'rose') {
      chart.setOption({
        title: { text: 'Holder Rose', left: 'center' },
        tooltip: { trigger: 'item' },
        legend: { top: 'bottom' },
        series: [{
          type: 'pie',
          radius: ['20%', '75%'],
          roseType: 'area',
          data: seriesData,
          label: { formatter: '{b}: {d}%'}
        }]
      });
    } else if (chartMode === 'treemap') {
      chart.setOption({
        title: { text: 'Holder Treemap', left: 'center' },
        tooltip: { formatter: '{b}: {c}' },
        series: [{
          type: 'treemap',
          data: seriesData,
          roam: false,
          label: { show: true, formatter: '{b}' }
        }]
      });
    } else if (chartMode === 'pareto') {
      const sorted = [...seriesData].sort((a, b) => b.value - a.value);
      const names = sorted.map(s => s.name);
      const values = sorted.map(s => s.value);
      const total = values.reduce((a, b) => a + b, 0) || 1;
      let acc = 0;
      const cum = values.map(v => {
        acc += v;
        return (acc / total * 100).toFixed(2);
      });
      chart.setOption({
        title: { text: 'Holder Pareto', left: 'center' },
        tooltip: { trigger: 'axis' },
        xAxis: [{ type: 'category', data: names, axisLabel: { rotate: 45 } }],
        yAxis: [{ type: 'value', name: 'Amount' }, { type: 'value', name: 'Cumulative %', min: 0, max: 100 }],
        series: [
          { type: 'bar', data: values, itemStyle: { color: '#2d8cf0' } },
          { type: 'line', yAxisIndex: 1, data: cum }
        ]
      });
    } else {
      chart.setOption({
        title: { text: 'Holder Distribution', left: 'center' },
        tooltip: { trigger: 'item' },
        legend: { top: 'bottom' },
        series: [{
          type: 'pie',
          radius: '60%',
          center: ['50%', '50%'],
          data: seriesData,
          label: { formatter: '{b}: {d}%'}
        }]
      });
    }
  </script>
</body>
</html>`
