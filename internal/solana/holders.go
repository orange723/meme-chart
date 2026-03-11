// Helius索引数据获取逻辑
package solana

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strconv"

	"meme-chart/internal/model"

	"github.com/tidwall/gjson"
)

// FetchTokenMeta 获取token元信息
func FetchTokenMeta(ctx context.Context, client *RPCClient, mint string) (model.TokenMeta, error) {
	// 调用getAsset
	resp, err := client.Call(ctx, "getAsset", map[string]interface{}{"id": mint})
	if err != nil {
		return model.TokenMeta{}, err
	}
	if errMsg := resp.Get("error.message"); errMsg.Exists() {
		return model.TokenMeta{}, fmt.Errorf("RPC错误: %s", errMsg.String())
	}

	// 解析名称与符号
	name := resp.Get("result.content.metadata.name").String()
	symbol := resp.Get("result.content.metadata.symbol").String()

	// 解析图片
	image := resp.Get("result.content.links.image").String()
	if image == "" {
		image = resp.Get("result.content.files.0.uri").String()
	}

	// 解析供应量与精度
	decimals := int(resp.Get("result.token_info.decimals").Int())
	supplyRaw := resp.Get("result.token_info.supply")
	supply := 0.0
	if supplyRaw.Exists() && supplyRaw.Type != gjson.Null {
		supply = supplyRaw.Float()
	}
	if supply == 0 {
		supplyStr := resp.Get("result.token_info.supply").String()
		if supplyStr != "" {
			if v, err := strconv.ParseFloat(supplyStr, 64); err == nil {
				supply = v
			}
		}
	}

	// 转换为UI数量
	if decimals > 0 {
		supply = supply / math.Pow10(decimals)
	}

	return model.TokenMeta{
		Name:     name,
		Symbol:   symbol,
		Mint:     mint,
		Image:    image,
		Supply:   supply,
		Decimals: decimals,
	}, nil
}

// FetchTopHolders 获取Top持仓数据(通过Helius索引)
func FetchTopHolders(ctx context.Context, client *RPCClient, mint string, top int, decimals int) ([]model.Holder, float64, error) {
	// 统计账户并聚合余额
	ownerAmounts := map[string]float64{}
	totalFromAccounts := 0.0
	cursor := ""

	for {
		// 组装请求
		params := map[string]interface{}{
			"mint":  mint,
			"limit": 1000,
			"options": map[string]interface{}{
				"showZeroBalance": false,
			},
		}
		if cursor != "" {
			params["cursor"] = cursor
		}

		// 调用getTokenAccounts
		resp, err := client.Call(ctx, "getTokenAccounts", params)
		if err != nil {
			return nil, 0, err
		}
		if errMsg := resp.Get("error.message"); errMsg.Exists() {
			return nil, 0, fmt.Errorf("RPC错误: %s", errMsg.String())
		}

		// 解析账户
		accounts := resp.Get("result.token_accounts").Array()
		if len(accounts) == 0 {
			break
		}
		for _, a := range accounts {
			owner := a.Get("owner").String()
			amountRaw := a.Get("amount")
			if owner == "" || !amountRaw.Exists() {
				continue
			}

			// amount为基础单位
			amount := amountRaw.Float()
			if decimals > 0 {
				amount = amount / math.Pow10(decimals)
			}
			if amount <= 0 {
				continue
			}
			ownerAmounts[owner] += amount
			totalFromAccounts += amount
		}

		// 读取cursor
		cursor = resp.Get("result.cursor").String()
		if cursor == "" {
			break
		}
	}

	// 排序取Top
	holders := sortHolders(ownerAmounts, top)
	return holders, totalFromAccounts, nil
}

// WithOthersBySupply 追加Others(优先使用token总量，若无则用账户合计)
func WithOthersBySupply(holders []model.Holder, totalSupply float64, totalFromAccounts float64) []model.Holder {
	// 计算TopN合计
	sum := 0.0
	for _, h := range holders {
		sum += h.UiAmount
	}

	// 计算Others
	total := totalSupply
	if total <= 0 {
		total = totalFromAccounts
	}
	others := total - sum
	if others < 0 {
		others = 0
	}
	return append(holders, model.Holder{Owner: "Others", TokenAccount: "", UiAmount: others})
}

// sortHolders 排序并截断TopN
func sortHolders(ownerAmounts map[string]float64, top int) []model.Holder {
	// 组装切片
	list := make([]model.Holder, 0, len(ownerAmounts))
	for owner, amount := range ownerAmounts {
		list = append(list, model.Holder{
			Owner:        shortenAddress(owner),
			TokenAccount: "",
			UiAmount:     amount,
		})
	}

	// 排序
	sort.Slice(list, func(i, j int) bool { return list[i].UiAmount > list[j].UiAmount })
	if top > 0 && len(list) > top {
		list = list[:top]
	}
	return list
}

// shortenAddress 缩短地址显示
func shortenAddress(addr string) string {
	// 地址太短就原样返回
	if len(addr) <= 10 {
		return addr
	}
	return addr[:4] + "..." + addr[len(addr)-4:]
}
