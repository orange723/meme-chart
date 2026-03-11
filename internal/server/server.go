// 动态HTTP服务
package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"

	"meme-chart/internal/model"
	"meme-chart/internal/solana"
)

// 页面模板数据
type pageData struct {
	Title     string
	Meta      model.TokenMeta
	Interval  int
	ChartMode string
}

// 页面模板
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
      <div class="sub">总量: {{printf "%.6f" .Meta.Supply}} | 合约: {{.Meta.Mint}} | 刷新: {{.Interval}} 秒</div>
    </div>
  </header>
  <div id="chart"></div>
  <script src="https://cdn.jsdelivr.net/npm/echarts@5/dist/echarts.min.js"></script>
  <script>
    const chart = echarts.init(document.getElementById('chart'));
    const chartMode = '{{.ChartMode}}';
    async function loadData() {
      const res = await fetch('/data');
      const data = await res.json();
      const holders = data.holders || [];
      const seriesData = holders.map(h => ({ name: h.owner, value: h.uiAmount }));
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
    }
    loadData();
    setInterval(loadData, {{.Interval}} * 1000);
  </script>
</body>
</html>`

// Start 启动动态服务
func Start(addr, endpoint, apiKey, mint string, top int, interval int, timeoutSec int, others bool, defaultChart string) error {
	// 初始化RPC客户端
	client, err := solana.NewRPCClient(endpoint, apiKey, time.Duration(timeoutSec)*time.Second)
	if err != nil {
		return err
	}

	// 初始化模板
	tmpl := template.Must(template.New("index").Parse(dynamicHTML))

	// 主页处理
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		meta, err := solana.FetchTokenMeta(ctx, client, mint)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		chartMode := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("chart")))
		if chartMode == "" {
			chartMode = strings.ToLower(strings.TrimSpace(defaultChart))
		}
		if chartMode == "" {
			chartMode = "pie"
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(w, pageData{Title: buildTitle(meta), Meta: meta, Interval: interval, ChartMode: chartMode}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// 数据接口
	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		meta, err := solana.FetchTokenMeta(ctx, client, mint)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		holders, totalFromAccounts, err := solana.FetchTopHolders(ctx, client, mint, top, meta.Decimals)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		if others {
			holders = solana.WithOthersBySupply(holders, meta.Supply, totalFromAccounts)
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"meta":      meta,
			"holders":   holders,
			"updatedAt": time.Now().Format(time.RFC3339),
		})
	})

	log.Printf("动态服务已启动: http://%s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		return fmt.Errorf("服务启动失败: %w", err)
	}
	return nil
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
