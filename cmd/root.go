// 命令行根命令
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// 支持的图表类型
var validChartModes = map[string]bool{
	"pie":     true,
	"bubble":  true,
	"donut":   true,
	"rose":    true,
	"treemap": true,
	"pareto":  true,
}

// 根命令对象
var rootCmd = &cobra.Command{
	Use:   "meme-chart",
	Short: "Solana meme token holders chart CLI",
	Long:  "获取Solana meme币持仓并渲染成HTML或提供动态服务",
}

// Execute 执行命令
func Execute() {
	// 执行并处理错误
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// validateChartMode 校验图表类型
func validateChartMode(mode string) error {
	if mode == "" {
		return nil // 允许为空，默认使用 pie
	}
	if !validChartModes[mode] {
		keys := make([]string, 0, len(validChartModes))
		for k := range validChartModes {
			keys = append(keys, k)
		}
		return fmt.Errorf("不支持的图表类型: %q，可选: %s", mode, strings.Join(keys, "/"))
	}
	return nil
}

// 初始化命令
func init() {
	// 设置环境变量前缀
	viper.SetEnvPrefix("MEMECHART")
	viper.AutomaticEnv()

	// 公共参数
	rootCmd.PersistentFlags().String("mint", "", "meme币Mint地址(必填)")
	rootCmd.PersistentFlags().String("endpoint", "https://mainnet.helius-rpc.com/", "Helius索引服务RPC地址")
	rootCmd.PersistentFlags().String("api-key", "", "Helius API Key(必填，或endpoint中已包含api-key)")
	rootCmd.PersistentFlags().Int("top", 20, "显示前多少个持仓")
	rootCmd.PersistentFlags().Int("timeout", 30, "RPC请求超时(秒)")
	rootCmd.PersistentFlags().Bool("others", false, "是否追加Others(总量减TopN)")
	rootCmd.PersistentFlags().String("chart", "pie", "图表类型(pie/bubble/donut/rose/treemap/pareto)")

	// 绑定viper
	bindFlags(map[string]string{
		"mint":     "mint",
		"endpoint": "endpoint",
		"api_key":  "api-key",
		"top":      "top",
		"timeout":  "timeout",
		"others":   "others",
		"chart":    "chart",
	})

	// 注册子命令
	rootCmd.AddCommand(renderCmd)
	rootCmd.AddCommand(serveCmd)
}

// bindFlags 绑定flag到viper，错误时输出警告
func bindFlags(pairs map[string]string) {
	for viperKey, flagName := range pairs {
		if err := viper.BindPFlag(viperKey, rootCmd.PersistentFlags().Lookup(flagName)); err != nil {
			fmt.Printf("警告: 绑定参数 %s 失败: %v\n", flagName, err)
		}
	}
}
