package http

import "github.com/gin-gonic/gin"

// SetupRoutes はGinルーターにすべてのAPIルートを設定するのです。
// ハンドラーの登録をここで一元管理するもなのです。
func SetupRoutes(router *gin.Engine) {
	// APIルートのグループ化
	api := router.Group("/api")
	{
		// ステータス取得
		api.GET("/status", GetStatus)

		// カレンダー取得
		api.GET("/calendar", GetCalendar)

		// タスク取得
		api.GET("/tasks", GetTasks)

		// 天気取得
		api.GET("/weather", GetWeather)

		// ヘルスチェック（疎通確認）
		api.GET("/health", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"ok": true,
			})
		})
	}

	// OAuth ルートのグループ化
	auth := router.Group("/auth")
	{
		// Google OAuth ログイン
		auth.GET("/login", AuthLogin)

		// Google OAuth コールバック
		auth.GET("/callback", AuthCallback)
	}
}
