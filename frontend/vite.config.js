import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

// https://vite.dev/config/
export default defineConfig({
  plugins: [svelte()],
  build: {
    // ビルド出力ディレクトリ（デフォルト dist から build に変更）
    outDir: 'build',
  },
  server: {
    // 開発時はバックエンドAPI（localhost:8080）にプロキシ
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/auth': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
  define: {
    // API ベース URL を環境変数から取得（デフォルト: バックエンド同じオリジン）
    'import.meta.env.VITE_API_BASE_URL': JSON.stringify(
      process.env.VITE_API_BASE_URL || 'http://localhost:8080'
    ),
  },
})
