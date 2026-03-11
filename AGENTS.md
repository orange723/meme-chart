# meme-chart

## 项目概览

基于 Golang 的 CLI 工具：通过 Helius 索引服务获取 Solana 上 meme 币持有分布与 token 元信息，渲染为静态 HTML 或启动动态服务。

## 目录结构

- `main.go`：程序入口
- `cmd/`：CLI 命令（cobra）
- `internal/solana/`：Helius API 调用与持有人数据处理
- `internal/render/`：静态 HTML 生成
- `internal/server/`：动态服务
- `internal/model/`：数据结构定义

## 依赖与工具

已使用依赖（保持列表真实、精简）：

- `github.com/spf13/cobra`
- `github.com/spf13/viper`
- `github.com/tidwall/gjson`

## 开发命令

- 运行静态渲染：`go run . render --mint <Mint地址> --api-key <HeliusApiKey> --top 10 --others --chart bubble`
- 运行动态服务：`go run . serve --mint <Mint地址> --api-key <HeliusApiKey> --top 10 --others --chart bubble`
- 格式化：`go fmt ./...`
- 静态检查：`go vet ./...`
- 测试：`go test ./...`

## 代码约定

- 每段代码，在上面写上中文注释，便于查看。
- 保持输出可预测：HTML 生成与服务响应应避免隐藏副作用。
- 新增依赖需同步更新 `go.mod` 与 `go.sum`。

## 外部服务与配置

- Helius API Key 必填：优先使用命令行参数 `--api-key`，也可用环境变量 `MEMECHART_API_KEY`。
- 访问网络仅用于调用 Helius API。
