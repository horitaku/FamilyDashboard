<script>
  import { onMount } from 'svelte';
  import { getCalendar } from '../api.js';

  let calendarData = null;
  let error = null;

  /**
   * カレンダーデータを取得
   */
  async function loadCalendarData() {
    try {
      calendarData = await getCalendar();
      error = null;
    } catch (err) {
      console.error('カレンダーデータ取得エラー:', err);
      error = err.message;
    }
  }

  /**
   * 日付をフォーマット（MM月DD日（曜日））
   */
  function formatDate(dateStr) {
    const date = new Date(dateStr + 'T00:00');
    const formatter = new Intl.DateTimeFormat('ja-JP', {
      month: 'numeric',
      day: 'numeric',
      weekday: 'short',
    });
    return formatter.format(date);
  }

  /**
   * 時刻をフォーマット（HH:MM）
   */
  function formatTime(dateStr) {
    if (!dateStr) return '';
    // ISO形式か確認して解析
    const date = new Date(dateStr);
    if (isNaN(date.getTime())) return '';
    const formatter = new Intl.DateTimeFormat('ja-JP', {
      hour: '2-digit',
      minute: '2-digit',
      hour12: false,
    });
    return formatter.format(date);
  }

  /**
   * イベントのスタイルを生成（色指定）
   */
  function getEventStyle(event) {
    let bgColor = '#e0e7ff'; // デフォルト：インジゴ
    let textColor = '#3730a3';

    // イベントの色情報がある場合は使用
    if (event.color) {
      bgColor = event.color;
      // テキスト色は自動判定（簡易版）
      const rgb = parseInt(event.color.slice(1), 16);
      const brightness = (rgb >> 16 & 255) * 0.299 + (rgb >> 8 & 255) * 0.587 + (rgb & 255) * 0.114;
      textColor = brightness > 128 ? '#000' : '#fff';
    }

    return `background-color: ${bgColor}; color: ${textColor}`;
  }

  onMount(() => {
    loadCalendarData();
    // 5分ごとにリロード
    const interval = setInterval(loadCalendarData, 300000);
    return () => clearInterval(interval);
  });
</script>

<div class="calendar-widget">
  {#if error}
    <div class="error">エラー: {error}</div>
  {:else if calendarData && calendarData.days}
    <div class="days-container">
      {#each calendarData.days as day}
        <div class="day">
          <div class="day-header">
            <h3 class="day-date">{formatDate(day.date)}</h3>
          </div>
          <div class="events-container">
            {#if day.allDay && day.allDay.length > 0}
              <div class="all-day-section">
                {#each day.allDay as event}
                  <div class="event all-day-event" style={getEventStyle(event)}>
                    <span class="event-title">{event.title}</span>
                  </div>
                {/each}
              </div>
            {/if}
            {#if day.timed && day.timed.length > 0}
              <div class="timed-events">
                {#each day.timed as event}
                  <div class="event timed-event" style={getEventStyle(event)}>
                    <span class="event-time">{formatTime(event.start)}</span>
                    <span class="event-title">{event.title}</span>
                  </div>
                {/each}
              </div>
            {/if}
            {#if (!day.allDay || day.allDay.length === 0) && (!day.timed || day.timed.length === 0)}
              <div class="no-events">予定なし</div>
            {/if}
          </div>
        </div>
      {/each}
    </div>
  {:else}
    <div class="loading">読み込み中...</div>
  {/if}
</div>

<style>
  .calendar-widget {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    padding: 24px;
    box-sizing: border-box;
    background: #fafafa;
    border-radius: 8px;
    overflow-y: auto;
  }

  .error,
  .loading {
    text-align: center;
    color: #666;
    font-size: 1.4rem;
    padding: 40px;
  }

  .error {
    color: #ef4444;
  }

  .days-container {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .day {
    border-left: 4px solid #3b82f6;
    padding-left: 16px;
  }

  .day-header {
    margin-bottom: 12px;
  }

  .day-date {
    font-size: 1.5rem;
    font-weight: bold;
    margin: 0;
    color: #1f2937;
  }

  .events-container {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .all-day-section {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .timed-events {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .event {
    padding: 12px 16px;
    border-radius: 4px;
    font-size: 1.1rem;
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .all-day-event {
    background: rgba(59, 130, 246, 0.15);
    border-left: 3px solid #3b82f6;
  }

  .timed-event {
    background: rgba(99, 102, 241, 0.1);
    border-left: 2px solid #6366f1;
  }

  .event-time {
    font-weight: bold;
    font-size: 1rem;
    min-width: 60px;
  }

  .event-title {
    flex: 1;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .no-events {
    color: #9ca3af;
    font-size: 1rem;
    font-style: italic;
  }
</style>
