# meme-chart

Solana meme 币持有分布可视化 CLI。基于 Helius 索引服务获取持有人分布与 token 元信息，生成静态 HTML 或启动动态服务。

## 功能

- Helius 索引数据：token 名称/符号/图片/总量
- 持有分布统计（Top N + Others）
- 多图表模式：`pie` / `bubble` / `donut` / `rose` / `treemap` / `pareto`
- 输出静态 HTML 或动态服务

## 安装与运行

### 1. 获取 Helius API Key

注册并创建 API Key：

```
https://dashboard.helius.dev
```

### 2. 运行（静态 HTML）

```bash
go run . render \
  --mint <Mint地址> \
  --api-key <HeliusApiKey> \
  --top 10 \
  --others \
  --chart bubble
```

输出文件名会自动使用 token 名称（若未指定 `--out`）。

### 3. 运行（动态服务）

```bash
go run . serve \
  --mint <Mint地址> \
  --api-key <HeliusApiKey> \
  --top 10 \
  --others \
  --chart bubble
```

访问：
```
http://127.0.0.1:8080
```

## 图表类型

使用 `--chart` 指定：

- `pie`：饼图
- `bubble`：泡泡图（类似 pump.fun）
- `donut`：环形图
- `rose`：玫瑰图
- `treemap`：矩形树图
- `pareto`：帕累托（条形 + 累计曲线）

## 版本发布

仓库包含 GitHub Actions Release 工作流：
- 当打 tag（如 `v0.1.0`）时自动编译多平台二进制并发布 Release

## 环境变量

你也可以用环境变量传递参数（前缀 `MEMECHART_`）：

```bash
export MEMECHART_API_KEY=xxxx
export MEMECHART_MINT=xxxx
```

## License

MIT
