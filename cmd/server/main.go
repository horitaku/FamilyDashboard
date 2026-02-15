package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rihow/FamilyDashboard/internal/cache"
	"github.com/rihow/FamilyDashboard/internal/config"
	httproutes "github.com/rihow/FamilyDashboard/internal/http"
	"github.com/rihow/FamilyDashboard/internal/services/google"
	"github.com/rihow/FamilyDashboard/internal/services/weather"
)

// main はGinサーバーのエントリーポイントなのです。
// 設定読み込み → APIルーティング → 静的ファイル配信 → サーバー起動 の順で処理するます。
func main() {
	// 設定ファイルを読み込むます。
	configFilePath := "./data/settings.json"
	cfg, err := config.LoadConfig(configFilePath)
	if err != nil {
		log.Fatalf("設定ファイルの読み込みに失敗しました: %v", err)
	}

	fmt.Printf("✨ 設定を読み込みました: %s\n", cfg.GetLocationString())
	fmt.Printf("   天気更新間隔: %v\n", cfg.GetRefreshInterval("weather"))
	fmt.Printf("   カレンダー更新間隔: %v\n", cfg.GetRefreshInterval("calendar"))
	fmt.Printf("   タスク更新間隔: %v\n", cfg.GetRefreshInterval("tasks"))

	// キャッシュを初期化するます
	fc := cache.New("./data/cache")

	// 天気APIクライアントを初期化するます
	weatherClient := weather.NewClient(fc, "http://localhost:8080")

	// Google APIクライアントを初期化するます
	googleClient := google.NewClient(fc, cfg)

	// 保存されたトークンを読み込む（以前にOAuth認可済みの場合）
	if err := googleClient.LoadTokens("./data/tokens.json"); err != nil {
		fmt.Printf("⚠️ トークン読込エラー: %v\n", err)
		// エラーでも継続する（トークンなしで開始してもOK）
	} else {
		fmt.Printf("✨ Google OAuth トークンを読み込みました\n")
	}

	// Ginルーターを初期化
	router := gin.Default()

	// グローバルミドルウェアで設定・クライアントをコンテキストに保存するます。
	router.Use(func(ctx *gin.Context) {
		ctx.Set("config", cfg)
		ctx.Set("weather", weatherClient)
		ctx.Set("google", googleClient)
		ctx.Next()
	})

	// APIルートの設定（internal/httpで定義したルートを登録）
	httproutes.SetupRoutes(router)
	router.Static("/assets", "./frontend/dist/assets")

	// ルートへのアクセスはindex.htmlを返す（SPA対応）
	router.NoRoute(func(ctx *gin.Context) {
		// indexファイルが存在しない場合は、エラーのみ返す（後でhosted filesになる予定）
		indexFile := "./frontend/dist/index.html"
		if _, err := os.Stat(indexFile); err == nil {
			ctx.File(indexFile)
		} else {
			ctx.JSON(404, gin.H{
				"error": "index.html not found. Frontend build required.",
			})
		}
	})

	// 既定ポート8080で起動するます。
	port := ":8080"
	fmt.Printf("🚀 サーバー起動するます！ http://localhost%s\n", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("サーバー起動に失敗しました: %v", err)
	}
}
