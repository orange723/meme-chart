# meme-chart

Solana meme 币持有分布可视化 CLI。基于 Helius 索引服务获取持有人分布与 token 元信息，生成静态 HTML 或启动动态服务。

## 功能

- **Helius 索引数据**：token 名称/符号/图片/总量
- **持有分布统计**：Top N + Others 聚合
- **6 种图表类型**：`pie` / `bubble` / `donut` / `rose` / `treemap` / `pareto`
- **双模式输出**：静态 HTML 文件 或 动态 HTTP 服务（自动刷新）
- **优雅退出**：动态服务支持 SIGINT/SIGTERM 平滑关闭
- **TokenMeta 缓存**：动态服务缓存元信息，减少重复请求

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

输出文件名自动使用 token 名称（可通过 `--out` 自定义）。

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

动态模式下图表每 30 秒自动刷新，可追加 `?chart=pie` 切换图表类型。

## 命令行参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--mint` | (必填) | meme 币 Mint 地址 |
| `--api-key` | (必填) | Helius API Key |
| `--endpoint` | `https://mainnet.helius-rpc.com/` | Helius RPC 地址 |
| `--top` | `20` | 显示前多少个持仓 |
| `--others` | `false` | 追加 Others（总量减 TopN） |
| `--chart` | `pie` | 图表类型 |
| `--timeout` | `30` | RPC 请求超时（秒） |
| `--out` | `meme-chart.html` | 渲染命令的输出文件 |
| `--addr` | `127.0.0.1:8080` | 服务命令的监听地址 |
| `--interval` | `30` | 服务命令的刷新间隔（秒） |

## 图表类型

使用 `--chart` 指定：

| 值 | 图表类型 |
|----|----------|
| `pie` | 饼图 |
| `bubble` | 泡泡图（力导向布局） |
| `donut` | 环形图 |
| `rose` | 玫瑰图（南丁格尔图） |
| `treemap` | 矩形树图 |
| `pareto` | 帕累托（条形 + 累计曲线） |

## 环境变量

所有参数均可通过环境变量传递（前缀 `MEMECHART_`）：

```bash
export MEMECHART_API_KEY=xxxx
export MEMECHART_MINT=xxxx
export MEMECHART_CHART=bubble
```

## 开发

```bash
# 构建
go build -o meme-chart .

# 格式化
go fmt ./...

# 静态检查
go vet ./...

# 测试
go test ./... -v -race

# 运行
go run . render --mint <地址> --api-key <key>
```

## CI/CD

- **CI**：push / PR 到 main 分支时自动运行 fmt / vet / test / build
- **Release**：打 tag（如 `v0.2.0`）时自动编译多平台二进制并发布 Release

## License

MIT
