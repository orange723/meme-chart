# Changelog

## [Unreleased]
- 消除模板重复：静态/动态模式共用 `render.DynamicHTML`
- 消除 `buildTitle` 函数重复，统一为 `render.BuildTitle`
- 动态服务缓存 TokenMeta，避免重复 RPC 请求
- 动态服务支持 SIGINT/SIGTERM 优雅退出
- 修复 `render` 命令 context 超时硬编码问题，统一使用 `--timeout`
- `--chart` 非法值时报错提示可用选项
- `viper.BindPFlag` 错误显式输出警告
- RPC 客户端使用原子自增 JSON-RPC ID
- `FetchTopHolders` 按需累加 `totalFromAccounts`（仅 `--others` 时）
- 提取 `defaultOutFile` 常量消除硬编码
- 新增单元测试（holders / client / render）
- 新增 CI workflow（push/PR 自动 fmt/vet/test/build）
- Release workflow 加入 vet/test 步骤

## [0.1.0] - 2026-03-11
- 初始可用版本
- Helius 索引获取 token 元信息与持仓
- 多图表渲染与 CLI 参数
