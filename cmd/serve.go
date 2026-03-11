// 动态服务命令
package cmd

import (
	"fmt"
	"strings"

	"meme-chart/internal/server"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// 动态服务命令对象
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "启动动态服务",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 读取参数
		mint := strings.TrimSpace(viper.GetString("mint"))
		endpoint := viper.GetString("endpoint")
		apiKey := viper.GetString("api_key")
		top := viper.GetInt("top")
		addr := viper.GetString("addr")
		interval := viper.GetInt("interval")
		timeoutSec := viper.GetInt("timeout")
		others := viper.GetBool("others")
		chartMode := strings.ToLower(strings.TrimSpace(viper.GetString("chart")))

		// 参数校验
		if mint == "" {
			return fmt.Errorf("必须提供 --mint 参数")
		}

		// 启动服务
		return server.Start(addr, endpoint, apiKey, mint, top, interval, timeoutSec, others, chartMode)
	},
}

// 初始化动态服务命令
func init() {
	// 专属参数
	serveCmd.Flags().String("addr", "127.0.0.1:8080", "服务监听地址")
	serveCmd.Flags().Int("interval", 30, "动态刷新间隔(秒)")

	// 绑定viper
	_ = viper.BindPFlag("addr", serveCmd.Flags().Lookup("addr"))
	_ = viper.BindPFlag("interval", serveCmd.Flags().Lookup("interval"))
}
