<script>
  import { onMount } from 'svelte';
  import { getStatus } from '../api.js';

  let now = new Date();
  let statusData = null;
  let hasError = false;

  /**
   * 現在時刻を Asia/Tokyo タイムゾーンで取得
   */
  function updateTime() {
    now = new Date();
  }

  /**
   * 時刻フォーマット（HH:MM）
   */
  function formatTime(date) {
    const formatter = new Intl.DateTimeFormat('ja-JP', {
      timeZone: 'Asia/Tokyo',
      hour: '2-digit',
      minute: '2-digit',
      hour12: false,
    });
    return formatter.format(date);
  }

  /**
   * 日付フォーマット（mm月dd日(曜日)）
   */
  function formatDate(date) {
    const formatter = new Intl.DateTimeFormat('ja-JP', {
      timeZone: 'Asia/Tokyo',
      month: 'numeric',
      day: 'numeric',
      weekday: 'short',
    });
    return formatter.format(date);
  }

  /**
   * ステータス情報を定期的に取得してエラー状態を監視
   */
  async function updateStatus() {
    try {
      statusData = await getStatus();
      // エラーがあれば hasError = true
      hasError = statusData.errors && statusData.errors.length > 0;
    } catch (error) {
      console.error('ステータス取得エラー:', error);
      hasError = true;
    }
  }

  onMount(() => {
    // 時刻を毎秒更新
    const timeInterval = setInterval(updateTime, 1000);

    // ステータスを初期ロード+5秒ごとに更新
    updateStatus();
    const statusInterval = setInterval(updateStatus, 5000);

    return () => {
      clearInterval(timeInterval);
      clearInterval(statusInterval);
    };
  });
</script>

<header>
  <div class="header-content">
    <div class="clock-container">
      <div class="time">{formatTime(now)}</div>
      <div class="date">{formatDate(now)}</div>
    </div>
    <div class="status-container" class:error={hasError}>
      {#if hasError}
        <div class="error-indicator"></div>
        <span class="error-text">接続エラー</span>
      {/if}
    </div>
  </div>
</header>

<style>
  header {
    width: 100%;
    height: 10%;
    background: linear-gradient(135deg, #1e3c72 0%, #2a5298 100%);
    display: flex;
    align-items: center;
    padding: 0 40px;
    box-sizing: border-box;
    color: white;
  }

  .header-content {
    width: 100%;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .clock-container {
    display: flex;
    flex-direction: column;
    justify-content: center;
  }

  .time {
    font-size: 5.5rem;
    font-weight: bold;
    line-height: 1;
    letter-spacing: 0.08em;
  }

  .date {
    font-size: 1.8rem;
    margin-top: 12px;
    opacity: 0.9;
    letter-spacing: 0.05em;
  }

  .status-container {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 16px 24px;
    border-radius: 8px;
    background: rgba(255, 255, 255, 0.1);
    min-width: 180px;
  }

  .status-container.error {
    background: rgba(239, 68, 68, 0.2);
  }

  .error-indicator {
    width: 16px;
    height: 16px;
    border-radius: 50%;
    background: #ef4444;
    animation: blink 2s ease-in-out infinite;
  }

  @keyframes blink {
    0%, 100% {
      opacity: 1;
    }
    50% {
      opacity: 0.3;
    }
  }

  .error-text {
    font-size: 1.2rem;
    color: #fca5a5;
  }

  /* エラーがない場合は非表示 */
  .status-container:not(.error) {
    visibility: hidden;
  }
</style>
