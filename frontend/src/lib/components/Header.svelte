<script>
  import { onMount } from 'svelte';
  import { getStatus } from '../api.js';

  let now = new Date();
  let statusData = null;
  let hasError = false;
  let isBackendOffline = false;
  let statusErrors = [];

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
   * ステータス情報を定期的に取得してエラー状態を監視
   */
  async function updateStatus() {
    try {
      statusData = await getStatus();
      isBackendOffline = false;
      statusErrors = Array.isArray(statusData.errors) ? statusData.errors : [];
      // エラーがあれば hasError = true
      hasError = statusErrors.length > 0;
    } catch (error) {
      console.error('ステータス取得エラー:', error);
      hasError = true;
      isBackendOffline = true;
      statusErrors = [];
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
    </div>
    <div class="status-container" class:error={hasError}>
      {#if hasError}
        <div class="error-indicator"></div>
        {#if isBackendOffline}
          <span class="error-text">サーバー接続エラー</span>
        {:else}
          <span class="error-text">データ取得エラー</span>
          {#if statusErrors.length > 0}
            <span class="error-sources">
              {statusErrors.map((item) => item.source).join(' / ')}
            </span>
          {/if}
        {/if}
      {/if}
    </div>
  </div>
</header>

<style>
  header {
    width: 100%;
    height: 5%;
    background: linear-gradient(135deg, #1e3c72 0%, #2a5298 100%);
    display: flex;
    align-items: center;
    padding: 0 16px;
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
    align-items: center;
  }

  .time {
    font-size: 2.5rem;
    font-weight: bold;
    line-height: 1;
    letter-spacing: 0.08em;
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
    animation: blink 2s steps(1, end) infinite;
  }

  @keyframes blink {
    0%, 49% {
      opacity: 1;
    }
    50%, 100% {
      opacity: 0;
    }
  }

  .error-text {
    font-size: 1.2rem;
    color: #fca5a5;
  }

  .error-sources {
    font-size: 0.95rem;
    color: #fecaca;
    letter-spacing: 0.02em;
  }

  /* エラーがない場合は非表示 */
  .status-container:not(.error) {
    visibility: hidden;
  }
</style>
