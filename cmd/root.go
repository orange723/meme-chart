// 命令行根命令
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
	_ = viper.BindPFlag("mint", rootCmd.PersistentFlags().Lookup("mint"))
	_ = viper.BindPFlag("endpoint", rootCmd.PersistentFlags().Lookup("endpoint"))
	_ = viper.BindPFlag("api_key", rootCmd.PersistentFlags().Lookup("api-key"))
	_ = viper.BindPFlag("top", rootCmd.PersistentFlags().Lookup("top"))
	_ = viper.BindPFlag("timeout", rootCmd.PersistentFlags().Lookup("timeout"))
	_ = viper.BindPFlag("others", rootCmd.PersistentFlags().Lookup("others"))
	_ = viper.BindPFlag("chart", rootCmd.PersistentFlags().Lookup("chart"))

	// 注册子命令
	rootCmd.AddCommand(renderCmd)
	rootCmd.AddCommand(serveCmd)
}
