/**
 * FamilyDashboard バックエンド APIクライアント
 * 
 * バックエンドサーバーの各エンドポイントにアクセスするためのクライアント
 * 全エンドポイントはバックエンドのみで提供される（フロントエンドは直接外部APIを呼ばない）
 */

// バックエンドのベースURL（デフォルト: 同じオリジン）
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

/**
 * APIエラークラス
 */
export class APIError extends Error {
  constructor(message, status, response) {
    super(message);
    this.status = status;
    this.response = response;
    this.name = 'APIError';
  }
}

/**
 * 汎用的なAPIリクエスト送信関数
 * @param {string} endpoint - エンドポイントパス（例: '/api/status'）
 * @param {object} options - fetchオプション
 * @returns {Promise<object>} レスポンスJSON
 */
async function request(endpoint, options = {}) {
  const url = `${API_BASE_URL}${endpoint}`;
  const defaultOptions = {
    headers: {
      'Content-Type': 'application/json',
    },
    ...options,
  };

  try {
    const response = await fetch(url, defaultOptions);

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new APIError(
        errorData.message || `HTTP ${response.status}`,
        response.status,
        errorData
      );
    }

    return await response.json();
  } catch (error) {
    if (error instanceof APIError) {
      throw error;
    }
    throw new APIError(`ネットワークエラー: ${error.message}`, 0, error);
  }
}

/**
 * ステータス情報を取得
 * @returns {Promise<object>} ステータス情報（ok, now, errors, lastUpdated）
 */
export async function getStatus() {
  return request('/api/status');
}

/**
 * カレンダーイベント/日程を取得
 * @returns {Promise<object>} カレンダーデータ（days配列）
 */
export async function getCalendar() {
  return request('/api/calendar');
}

/**
 * タスク一覧を取得
 * @returns {Promise<object>} タスクデータ（items配列）
 */
export async function getTasks() {
  return request('/api/tasks');
}

/**
 * 天気情報を取得
 * @returns {Promise<object>} 天気データ（current, today, alerts等）
 */
export async function getWeather() {
  return request('/api/weather');
}
