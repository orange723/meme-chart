# meme-chart

## 项目描述

以 golang 为编程语言的 CLI 工具，使用 Helius 索引服务获取 Solana 上 meme 币持有分布，并渲染为 HTML（静态或动态）。支持多种图表（饼图/泡泡图/玫瑰图/矩形树图/帕累托）。

## 项目实现

- cli 项目（golang + cobra/viper）
- 通过 Helius 索引 API 获取：
  - token 元信息（名称/符号/图片/总量）
  - 持有人分布
- 输出静态 HTML 或启动动态服务

## 工具

可使用以下包（已使用）：

- https://github.com/go-echarts/go-echarts/tree/master
- https://github.com/imroc/req
- https://github.com/spf13/viper
- https://github.com/spf13/cobra
- https://github.com/tidwall/gjson

## 协作偏好

- 每段代码，在上面写上中文注释，便于查看
