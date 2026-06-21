// 静态渲染命令
package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"meme-chart/internal/model"
	"meme-chart/internal/render"
	"meme-chart/internal/solana"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// defaultOutFile 默认输出文件名
const defaultOutFile = "meme-chart.html"

// 静态渲染命令对象
var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "渲染静态HTML图表",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 读取参数
		mint := strings.TrimSpace(viper.GetString("mint"))
		endpoint := viper.GetString("endpoint")
		apiKey := viper.GetString("api_key")
		top := viper.GetInt("top")
		out := viper.GetString("out")
		timeoutSec := viper.GetInt("timeout")
		others := viper.GetBool("others")
		chartMode := strings.ToLower(strings.TrimSpace(viper.GetString("chart")))

		// 参数校验
		if mint == "" {
			return fmt.Errorf("必须提供 --mint 参数")
		}
		if err := validateChartMode(chartMode); err != nil {
			return err
		}

		// 确保最小超时
		if timeoutSec < 5 {
			timeoutSec = 5
		}
		timeout := time.Duration(timeoutSec) * time.Second

		// 获取数据（使用用户指定的超时）
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		client, err := solana.NewRPCClient(endpoint, apiKey, timeout)
		if err != nil {
			return err
		}
		meta, err := solana.FetchTokenMeta(ctx, client, mint)
		if err != nil {
			return err
		}
		holders, totalFromAccounts, err := solana.FetchTopHolders(ctx, client, mint, top, meta.Decimals, others)
		if err != nil {
			return err
		}
		if others {
			holders = solana.WithOthersBySupply(holders, meta.Supply, totalFromAccounts)
		}

		// 生成默认输出文件名
		if out == defaultOutFile {
			out = buildOutputFilename(meta)
		}

		// 渲染输出
		if err := render.RenderStaticHTML(holders, meta, chartMode, out); err != nil {
			return err
		}

		fmt.Printf("已生成: %s\n", out)
		return nil
	},
}

// 初始化静态渲染命令
func init() {
	// 专属参数
	renderCmd.Flags().String("out", defaultOutFile, "输出HTML文件")

	// 绑定viper
	if err := viper.BindPFlag("out", renderCmd.Flags().Lookup("out")); err != nil {
		fmt.Printf("警告: 绑定 out 参数失败: %v\n", err)
	}
}

// buildOutputFilename 构建输出文件名
func buildOutputFilename(meta model.TokenMeta) string {
	// 选择基础名称
	base := meta.Name
	if base == "" {
		base = meta.Symbol
	}
	if base == "" {
		base = meta.Mint
	}

	// 生成安全文件名
	safe := sanitizeFilename(base)
	if safe == "" {
		safe = "token"
	}
	return safe + ".html"
}

// sanitizeFilename 生成安全文件名(ASCII)
func sanitizeFilename(s string) string {
	// 统一小写，空格转-
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")

	// 过滤非ASCII字符
	var b strings.Builder
	for _, r := range s {
		if r >= 'a' && r <= 'z' {
			b.WriteRune(r)
			continue
		}
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
			continue
		}
		if r == '-' || r == '_' || r == '.' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
