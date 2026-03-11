// 数据模型定义
package model

// Holder 持仓数据
type Holder struct {
	Owner        string  `json:"owner"`
	TokenAccount string  `json:"tokenAccount"`
	UiAmount     float64 `json:"uiAmount"`
}

// TokenMeta token元信息
type TokenMeta struct {
	Name     string  `json:"name"`
	Symbol   string  `json:"symbol"`
	Mint     string  `json:"mint"`
	Image    string  `json:"image"`
	Supply   float64 `json:"supply"`
	Decimals int     `json:"decimals"`
}
