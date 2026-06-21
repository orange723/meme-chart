// 共享页面模板与工具函数
package render

import (
	"html/template"

	"meme-chart/internal/model"
)

// PageData 页面模板数据
// Interval=0 为静态模式（内嵌数据），Interval>0 为动态模式（AJAX 刷新）
type PageData struct {
	Title       string
	Meta        model.TokenMeta
	ChartMode   string
	Interval    int         // 动态模式刷新间隔（秒），0=静态
	PayloadJSON template.JS // 静态模式内嵌的JSON数据
}

// BuildTitle 构建页面标题
func BuildTitle(meta model.TokenMeta) string {
	// 优先使用名称+符号
	if meta.Name != "" && meta.Symbol != "" {
		return meta.Name + " (" + meta.Symbol + ")"
	}
	if meta.Name != "" {
		return meta.Name
	}
	return meta.Mint
}

// 统一页面模板（静态与动态共用）
const DynamicHTML = `<!DOCTYPE html>
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
      <div class="sub">总量: {{printf "%.6f" .Meta.Supply}} | 合约: {{.Meta.Mint}}{{if .Interval}} | 刷新: {{.Interval}} 秒{{end}}</div>
    </div>
  </header>
  <div id="chart"></div>
  <script src="https://cdn.jsdelivr.net/npm/echarts@5/dist/echarts.min.js"></script>
  <script>
    const chart = echarts.init(document.getElementById('chart'));
    const defaultChartMode = '{{.ChartMode}}';
    const palette = ['#1f77b4','#ff7f0e','#2ca02c','#d62728','#9467bd','#8c564b','#e377c2','#7f7f7f','#bcbd22','#17becf'];

    // 渲染图表
    function renderChart(holders, chartMode) {
      const seriesData = holders.map(h => ({ name: h.owner, value: h.uiAmount }));
      if (chartMode === 'bubble') {
        const max = Math.max(...seriesData.map(s => s.value), 1);
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
        // 默认 pie
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
    }

    {{if .Interval}}
    // 动态模式：AJAX 轮询
    async function loadData() {
      const urlParams = new URLSearchParams(window.location.search);
      const chartParam = urlParams.get('chart');
      const chartMode = chartParam || defaultChartMode;
      try {
        const res = await fetch('/data');
        if (!res.ok) throw new Error('HTTP ' + res.status);
        const data = await res.json();
        renderChart(data.holders || [], chartMode);
      } catch (e) {
        console.error('数据加载失败:', e);
      }
    }
    loadData();
    setInterval(loadData, {{.Interval}} * 1000);
    {{else}}
    // 静态模式：内嵌数据
    const payload = {{.PayloadJSON}};
    const chartMode = payload.chartMode || defaultChartMode;
    renderChart(payload.holders || [], chartMode);
    {{end}}
  </script>
</body>
</html>`
