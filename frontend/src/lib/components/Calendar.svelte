<script>
  import { onMount } from 'svelte';
  import { getCalendar } from '../api.js';

  // Props: 表示する日数、スキップする日数、タイトル
  export let daysToShow = 7; // デフォルト: 7日分
  export let skipDays = 0;   // デフォルト: 0日スキップ
  export let title = '';     // デフォルト: タイトルなし
  export let showDate = false; // デフォルト: 日付非表示

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
   * 表示する日数を絞り込む
   */
  function filterDays(days) {
    if (!days) return [];
    return days.slice(skipDays, skipDays + daysToShow);
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
   * 簡易日付フォーマット (M/D)
   */
  function formatShortDate(dateString) {
    const date = new Date(dateString);
    const month = date.getMonth() + 1;
    const day = date.getDate();
    return `${month}/${day}`;
  }

  /**
   * イベントのスタイルを生成（色を枠線に使用）
   */
  function getEventStyle(event) {
    let borderColor = '#3b82f6'; // デフォルト：青

    // イベントの色情報がある場合は枠線に使用
    if (event.color) {
      borderColor = event.color;
    }

    return `border-left-color: ${borderColor}; background-color: #ffffff; color: #1f2937;`;
  }

  onMount(() => {
    loadCalendarData();
    // 5分ごとにリロード
    const interval = setInterval(loadCalendarData, 300000);
    return () => clearInterval(interval);
  });
</script>

<div class="calendar-widget">
  {#if title}
    <h2 class="widget-title">{title}</h2>
  {/if}
  {#if error}
    <div class="error">エラー: {error}</div>
  {:else if calendarData && calendarData.days}
    <div class="days-container">
      {#each filterDays(calendarData.days) as day}
        <div class="day">
          {#if !title}
            <div class="day-header">
              <h3 class="day-date">{formatDate(day.date)}</h3>
            </div>
          {/if}
          <div class="events-container">
            {#if day.allDay && day.allDay.length > 0}
              <div class="all-day-section">
                {#each day.allDay as event}
                  <div class="event all-day-event" style={getEventStyle(event)}>
                    {#if showDate}
                      <span class="event-date">{formatShortDate(day.date)}</span>
                    {/if}
                    <span class="event-title">{event.title}</span>
                  </div>
                {/each}
              </div>
            {/if}
            {#if day.timed && day.timed.length > 0}
              <div class="timed-events">
                {#each day.timed as event}
                  <div class="event timed-event" style={getEventStyle(event)}>
                    {#if showDate}
                      <span class="event-date">{formatShortDate(day.date)}</span>
                    {/if}
                    <span class="event-time">
                      <span class="time-start">{formatTime(event.start)}</span>
                      <span class="time-end">{formatTime(event.end)}</span>
                    </span>
                    <span class="event-title">{event.title}</span>
                  </div>
                {/each}
              </div>
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
    padding: 8px 24px 24px 24px;
    box-sizing: border-box;
    background: linear-gradient(135deg, #f0f3f8 0%, #e0e6f0 100%);
    border-radius: 8px;
    overflow-y: auto;
  }

  .widget-title {
    margin: 0 0 8px 0;
    font-size: 1.6rem;
    font-weight: bold;
    color: #1f2937;
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
    padding: 8px 10px;
    border-radius: 4px;
    font-size: 1.3rem;
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .all-day-event {
    background: #ffffff;
    border-left: 8px solid #3b82f6;
    color: #1f2937;
  }

  .timed-event {
    background: #ffffff;
    border-left: 8px solid #6366f1;
    color: #1f2937;
  }

  .event-date {
    font-weight: bold;
    font-size: 0.9rem;
    color: #6b7280;
    white-space: nowrap;
    min-width: 40px;
    text-align: center;
  }

  .event-time {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    font-weight: bold;
    font-size: 0.75rem;
    min-width: 38px;
    line-height: 1.2;
    gap: 2px;
  }

  .event-title {
    flex: 1;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
</style>
