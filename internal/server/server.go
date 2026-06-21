// 动态HTTP服务
package server

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"meme-chart/internal/model"
	"meme-chart/internal/render"
	"meme-chart/internal/solana"
)

// Start 启动动态服务
func Start(addr, endpoint, apiKey, mint string, top int, interval int, timeoutSec int, others bool, defaultChart string) error {
	// 初始化RPC客户端
	client, err := solana.NewRPCClient(endpoint, apiKey, time.Duration(timeoutSec)*time.Second)
	if err != nil {
		return err
	}

	// TokenMeta 缓存（启动时请求一次，后续复用）
	var (
		metaCache    model.TokenMeta
		metaOnce     sync.Once
		metaFetchErr error
	)

	// 获取缓存的元信息
	getCachedMeta := func(ctx context.Context) (model.TokenMeta, error) {
		metaOnce.Do(func() {
			metaCache, metaFetchErr = solana.FetchTokenMeta(ctx, client, mint)
		})
		return metaCache, metaFetchErr
	}

	// 初始化模板
	tmpl := template.Must(template.New("index").Parse(render.DynamicHTML))

	// 主页处理
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		meta, err := getCachedMeta(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		chartMode := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("chart")))
		if chartMode == "" {
			chartMode = strings.ToLower(strings.TrimSpace(defaultChart))
		}
		if chartMode == "" {
			chartMode = "pie"
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(w, render.PageData{
			Title:     render.BuildTitle(meta),
			Meta:      meta,
			ChartMode: chartMode,
			Interval:  interval,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// 数据接口
	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		meta, err := getCachedMeta(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		holders, totalFromAccounts, err := solana.FetchTopHolders(ctx, client, mint, top, meta.Decimals, others)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		if others {
			holders = solana.WithOthersBySupply(holders, meta.Supply, totalFromAccounts)
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"meta":      meta,
			"holders":   holders,
			"updatedAt": time.Now().Format(time.RFC3339),
		})
	})

	// 启动HTTP服务（支持优雅退出）
	srv := &http.Server{Addr: addr}

	// 监听退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("正在关闭服务...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("关闭异常: %v", err)
		}
	}()

	log.Printf("动态服务已启动: http://%s", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("服务启动失败: %w", err)
	}
	log.Println("服务已停止")
	return nil
}
