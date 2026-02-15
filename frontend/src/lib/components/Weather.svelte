<script>
  import { onMount } from 'svelte';
  import { getWeather } from '../api.js';

  let weatherData = null;
  let error = null;

  /**
   * 天気データを取得
   */
  async function loadWeatherData() {
    try {
      weatherData = await getWeather();
      error = null;
    } catch (err) {
      console.error('天気データ取得エラー:', err);
      error = err.message;
    }
  }

  /**
   * ウェザーアイコンを絵文字で表示
   */
  function getWeatherEmoji(iconCode) {
    // WMO天気コードとアイコン形式（01d など）で判定
    if (!iconCode) return '🌡️';
    
    const code = iconCode.substring(0, 2);
    switch (code) {
      case '01':
        return '☀️'; // 晴・clear
      case '02':
        return '🌤️'; // 薄曇・partly cloudy
      case '03':
      case '04':
        return '☁️'; // 曇り・cloudy
      case '09':
      case '10':
        return '🌧️'; // 雨・rain
      case '11':
        return '⛈️'; // 雷雨・thunderstorm
      case '13':
        return '❄️'; // 雪・snow
      case '50':
        return '🌫️'; // 霧・mist
      default:
        return '🌡️';
    }
  }

  onMount(() => {
    loadWeatherData();
    // 5分ごとにリロード
    const interval = setInterval(loadWeatherData, 300000);
    return () => clearInterval(interval);
  });
</script>

<div class="weather-widget">
  {#if error}
    <div class="error">エラー: {error}</div>
  {:else if weatherData}
    <!-- 現在の天候セクション -->
    <div class="current-section">
      <div class="current-icon">{getWeatherEmoji(weatherData.current?.icon || '')}</div>
      <div class="current-info">
        <div class="current-temp">{Math.round(weatherData.current?.temperature || 0)}°C</div>
        <div class="current-condition">{weatherData.current?.condition || '---'}</div>
      </div>
    </div>

    <!-- 本日の最高・最低気温 -->
    {#if weatherData.today}
      <div class="today-section">
        <div class="temp-row">
          <span class="temp-label">最高</span>
          <span class="temp-value">{Math.round(weatherData.today.maxTemp || 0)}°C</span>
        </div>
        <div class="temp-row">
          <span class="temp-label">最低</span>
          <span class="temp-value">{Math.round(weatherData.today.minTemp || 0)}°C</span>
        </div>
      </div>
    {/if}

    <!-- 降水確率 -->
    {#if weatherData.precipSlots && weatherData.precipSlots.length > 0}
      <div class="precip-section">
        <h4 class="section-title">降水確率</h4>
        <div class="precip-grid">
          {#each weatherData.precipSlots.slice(0, 8) as slot}
            <div class="precip-slot">
              <div class="precip-time">{slot.time}</div>
              <div class="precip-bar-container">
                <div
                  class="precip-bar"
                  style={`height: ${slot.probability}%`}
                ></div>
              </div>
              <div class="precip-percent">{slot.probability}%</div>
            </div>
          {/each}
        </div>
      </div>
    {/if}

    <!-- アラート・注意報 -->
    {#if weatherData.alerts && weatherData.alerts.length > 0}
      <div class="alerts-section">
        {#each weatherData.alerts as alert}
          <div class="alert alert-{alert.severity}">
            <span class="alert-severity">{alert.severity}</span>
            <span class="alert-title">{alert.title}</span>
          </div>
        {/each}
      </div>
    {/if}
  {:else}
    <div class="loading">読み込み中...</div>
  {/if}
</div>

<style>
  .weather-widget {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    padding: 24px;
    box-sizing: border-box;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    border-radius: 8px;
    overflow-y: auto;
  }

  .error,
  .loading {
    text-align: center;
    color: rgba(255, 255, 255, 0.8);
    font-size: 1.2rem;
    padding: 40px;
  }

  .error {
    color: #fca5a5;
  }

  /* 現在の天候 */
  .current-section {
    display: flex;
    align-items: center;
    gap: 20px;
    margin-bottom: 24px;
    padding-bottom: 20px;
    border-bottom: 2px solid rgba(255, 255, 255, 0.2);
  }

  .current-icon {
    font-size: 5rem;
    line-height: 1;
  }

  .current-info {
    display: flex;
    flex-direction: column;
    justify-content: center;
  }

  .current-temp {
    font-size: 3rem;
    font-weight: bold;
    line-height: 1;
  }

  .current-condition {
    font-size: 1.3rem;
    margin-top: 8px;
    opacity: 0.95;
  }

  /* 最高・最低気温 */
  .today-section {
    display: flex;
    gap: 24px;
    margin-bottom: 20px;
    padding: 16px;
    background: rgba(255, 255, 255, 0.1);
    border-radius: 6px;
  }

  .temp-row {
    flex: 1;
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 1.2rem;
  }

  .temp-label {
    opacity: 0.8;
  }

  .temp-value {
    font-weight: bold;
    font-size: 1.5rem;
  }

  /* 降水確率 */
  .precip-section {
    margin-bottom: 20px;
  }

  .section-title {
    margin: 0 0 12px 0;
    font-size: 1rem;
    opacity: 0.9;
  }

  .precip-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 8px;
  }

  .precip-slot {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 6px;
  }

  .precip-time {
    font-size: 0.85rem;
    opacity: 0.8;
  }

  .precip-bar-container {
    width: 24px;
    height: 80px;
    background: rgba(255, 255, 255, 0.2);
    border-radius: 3px;
    display: flex;
    align-items: flex-end;
    overflow: hidden;
  }

  .precip-bar {
    width: 100%;
    background: #fbbf24;
    border-radius: 2px;
    transition: height 0.3s ease;
  }

  .precip-percent {
    font-size: 0.8rem;
    font-weight: bold;
  }

  /* アラート */
  .alerts-section {
    display: flex;
    flex-direction: column;
    gap: 8px;
    padding-top: 16px;
    border-top: 2px solid rgba(255, 255, 255, 0.2);
  }

  .alert {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px;
    border-radius: 4px;
    font-size: 0.95rem;
    animation: blink-alert 2s ease-in-out infinite;
  }

  .alert-warning {
    background: rgba(245, 158, 11, 0.3);
    border-left: 3px solid #fbbf24;
  }

  .alert-error {
    background: rgba(239, 68, 68, 0.3);
    border-left: 3px solid #fca5a5;
  }

  .alert-severity {
    font-weight: bold;
    font-size: 0.85rem;
    text-transform: uppercase;
    min-width: 50px;
  }

  .alert-title {
    flex: 1;
  }

  @keyframes blink-alert {
    0%, 100% {
      opacity: 1;
    }
    50% {
      opacity: 0.6;
    }
  }
</style>
