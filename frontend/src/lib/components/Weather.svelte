  <script>
  import { onMount } from 'svelte';
  import { getWeather } from '../api.js';

  const WEATHER_ICON_VERSION = '20260307';

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
   * アイコンコードをローカルSVGパスに変換
   */
  function withIconVersion(path) {
    return `${path}?v=${WEATHER_ICON_VERSION}`;
  }

  function getWeatherIconPath(iconCode) {
    if (!iconCode) return withIconVersion('/weather-icons/unknown.svg');

    const code = iconCode.substring(0, 2);
    switch (code) {
      case '01':
        return withIconVersion('/weather-icons/clear.svg');
      case '02':
        return withIconVersion('/weather-icons/partly-cloudy.svg');
      case '03':
        return withIconVersion('/weather-icons/cloudy.svg');
      case '09':
        return withIconVersion('/weather-icons/drizzle.svg');
      case '10':
        return withIconVersion('/weather-icons/rain.svg');
      case '11':
        return withIconVersion('/weather-icons/heavy-rain.svg');
      case '12':
        return withIconVersion('/weather-icons/shower.svg');
      case '13':
        return withIconVersion('/weather-icons/snow.svg');
      case '14':
        return withIconVersion('/weather-icons/blizzard.svg');
      case '15':
        return withIconVersion('/weather-icons/thunder.svg');
      case '50':
        return withIconVersion('/weather-icons/fog.svg');
      default:
        return withIconVersion('/weather-icons/unknown.svg');
    }
  }

  /**
   * 週次用の曜日表示
   */
  function formatWeekday(dateStr) {
    if (!dateStr) return '';
    const date = new Date(`${dateStr}T00:00:00+09:00`);
    const formatter = new Intl.DateTimeFormat('ja-JP', {
      timeZone: 'Asia/Tokyo',
      weekday: 'short',
    });
    return formatter.format(date);
  }

  /**
   * 降水確率の値を取得
   */
  function getPrecipValue(slot) {
    if (!slot) return 0;
    if (typeof slot.probability === 'number') return slot.probability;
    if (typeof slot.precip === 'number') return slot.precip;
    return 0;
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
    <div class="current-block">
      <div class="current-main">
        <div class="current-icon">
          <img
            src={getWeatherIconPath(weatherData.current?.icon || '')}
            alt={weatherData.current?.condition || '天気アイコン'}
            loading="lazy"
          />
        </div>
        <div class="current-text">
          <div class="current-condition">{weatherData.current?.condition || '---'}</div>
          <div class="current-location">{weatherData.location || ''}</div>
        </div>
      </div>
      {#if weatherData.today}
        <div class="current-temps">
          <div class="temp-card temp-now">
            <span class="temp-label">いまのきおん</span>
            <span class="temp-value">{Math.round(weatherData.current?.temperature || 0)}°C</span>
          </div>
          <div class="temp-card temp-max">
            <span class="temp-label">さいこう</span>
            <span class="temp-value">{Math.round(weatherData.today.maxTemp || 0)}°C</span>
          </div>
          <div class="temp-card temp-min">
            <span class="temp-label">さいてい</span>
            <span class="temp-value">{Math.round(weatherData.today.minTemp || 0)}°C</span>
          </div>
        </div>
      {/if}
    </div>

    <div class="alerts-section">
      {#if weatherData.alerts && weatherData.alerts.length > 0}
        {#each weatherData.alerts as alert}
          <div class="alert {alert.severity === '特別警報' ? 'alert-special' : alert.severity === '警報' ? 'alert-error' : 'alert-warning'}">
            <span class="alert-severity">{alert.severity}</span>
            <span class="alert-title">{alert.title}</span>
          </div>
        {/each}
      {:else}
        <div class="alert alert-ok">
          現在、注意報・警報はありません
        </div>
      {/if}
    </div>

    {#if weatherData.precipSlots && weatherData.precipSlots.length > 0}
      <div class="hourly-section">
        <div class="hourly-grid">
          {#each weatherData.precipSlots.slice(0, 8) as slot}
            <div class="hourly-slot">
              <div class="hourly-time">{slot.time}</div>
              <div class="hourly-icon">
                <img
                  src={getWeatherIconPath(slot.icon || weatherData.current?.icon || '')}
                  alt="時間帯の天気アイコン"
                  loading="lazy"
                />
              </div>
              <div class="hourly-precip">{getPrecipValue(slot)}%</div>
            </div>
          {/each}
        </div>
      </div>
    {/if}
    
    {#if weatherData.weekly && weatherData.weekly.length > 0}
      <div class="weekly-section">
        <div class="weekly-grid">
          {#each weatherData.weekly.slice(0, 7) as day}
            <div class="weekly-item">
              <div class="weekly-day">{formatWeekday(day.date)}</div>
              <div class="weekly-icon">
                <img
                  src={getWeatherIconPath(day.icon || '')}
                  alt="週間天気アイコン"
                  loading="lazy"
                />
              </div>
              <div class="weekly-temps">
                <span class="weekly-max">{Math.round(day.maxTemp || 0)}°</span>
                <span class="weekly-min">{Math.round(day.minTemp || 0)}°</span>
              </div>
            </div>
          {/each}
        </div>
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
    gap: 8px;
    padding: 8px;
    box-sizing: border-box;
    background: linear-gradient(135deg, #0f3d66 0%, #1d6fa5 52%, #46b6c9 100%);
    color: #f8fafc;
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

  .current-block {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 16px;
    padding: 0px 12px;
    background: rgba(255, 255, 255, 0.12);
    border-radius: 12px;
    border: 1px solid rgba(255, 255, 255, 0.14);
  }

  .current-main {
    display: flex;
    align-items: center;
    gap: 16px;
    flex: 1;
    min-width: 0;
  }

  .current-icon {
    width: 7.8rem;
    height: 7.8rem;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .current-icon img {
    width: 100%;
    height: 100%;
    object-fit: contain;
  }

  .current-text {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .current-condition {
    font-size: 2.2rem;
    font-weight: 700;
    letter-spacing: 0.02em;
  }

  .current-location {
    font-size: 1rem;
    opacity: 0.85;
  }

  .current-temps {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 10px;
  }

  .temp-card {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 8px 12px;
    border-radius: 10px;
    min-width: 86px;
  }

  .temp-now {
    background: rgba(14, 165, 233, 0.22);
    border: 1px solid rgba(125, 211, 252, 0.4);
  }

  .temp-max {
    background: rgba(253, 230, 138, 0.22);
    border: 1px solid rgba(253, 230, 138, 0.4);
  }

  .temp-min {
    background: rgba(191, 219, 254, 0.18);
    border: 1px solid rgba(191, 219, 254, 0.4);
  }

  .temp-label {
    font-size: 0.85rem;
    opacity: 0.9;
    letter-spacing: 0.08em;
  }

  .temp-value {
    font-size: 2rem;
    font-weight: 700;
  }

  .hourly-section {
    padding: 12px;
    background: rgba(255, 255, 255, 0.08);
    border-radius: 12px;
  }

  .hourly-grid {
    display: grid;
    grid-template-columns: repeat(8, 1fr);
    gap: 10px;
  }

  .hourly-slot {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 8px 4px;
    background: rgba(15, 23, 42, 0.2);
    border-radius: 10px;
  }

  .hourly-time {
    font-size: 1.275rem;
    opacity: 0.85;
    font-weight: 600;
  }

  .hourly-icon {
    width: 3.6rem;
    height: 3.6rem;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .hourly-icon img {
    width: 100%;
    height: 100%;
    object-fit: contain;
  }

  .hourly-precip {
    font-size: 1.35rem;
    font-weight: 700;
  }

  .weekly-section {
    padding: 12px;
    background: rgba(255, 255, 255, 0.08);
    border-radius: 12px;
  }

  .weekly-grid {
    display: grid;
    grid-template-columns: repeat(7, 1fr);
    gap: 8px;
  }

  .weekly-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 6px 0;
    background: rgba(15, 23, 42, 0.18);
    border-radius: 8px;
  }

  .weekly-day {
    font-size: 1.2rem;
    opacity: 0.85;
    font-weight: 600;
  }

  .weekly-icon {
    width: 2.7rem;
    height: 2.7rem;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .weekly-icon img {
    width: 100%;
    height: 100%;
    object-fit: contain;
  }

  .weekly-temps {
    display: flex;
    gap: 4px;
    font-size: 1.125rem;
    font-weight: 700;
  }

  .weekly-max {
    color: #fde68a;
  }

  .weekly-min {
    color: #bfdbfe;
  }

  /* アラート */
  .alerts-section {
    display: flex;
    flex-direction: column;
    gap: 8px;
    padding-top: 0;
  }

  .alert {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px;
    border-radius: 8px;
    font-size: 1.425rem;
    animation: blink-alert 2s ease-in-out infinite;
  }

  .alert-warning {
    background: rgba(245, 158, 11, 0.7);
    border-left: 6px solid #fbbf24;
  }

  .alert-special {
    background: rgba(168, 85, 247, 0.7);
    border-left: 6px solid #d8b4fe;
  }

  .alert-error {
    background: rgba(239, 68, 68, 0.7);
    border-left: 6px solid #fca5a5;
  }

  .alert-ok {
    background: rgba(16, 185, 129, 0.2);
    border-left: 6px solid #6ee7b7;
    font-size: 1.35rem;
    animation: none;
  }

  .alert-severity {
    font-weight: bold;
    font-size: 1.275rem;
    text-transform: uppercase;
    min-width: 75px;
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
