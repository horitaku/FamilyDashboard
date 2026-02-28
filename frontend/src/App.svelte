<script>
  import Header from './lib/components/Header.svelte'
  import Calendar from './lib/components/Calendar.svelte'
  import Weather from './lib/components/Weather.svelte'
  import Tasks from './lib/components/Tasks.svelte'

  /**
   * 今日の日付を「M/D(曜)」形式でフォーマット
   */
  function getTodayDateString() {
    const today = new Date();
    const month = today.getMonth() + 1;
    const day = today.getDate();
    const weekdayFormatter = new Intl.DateTimeFormat('ja-JP', {
      timeZone: 'Asia/Tokyo',
      weekday: 'short',
    });
    const weekday = weekdayFormatter.format(today);
    return `${month}/${day}(${weekday})`;
  }

  /**
   * 明日の日付を「M/D(曜)」形式でフォーマット
   */
  function getTomorrowDateString() {
    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);
    const month = tomorrow.getMonth() + 1;
    const day = tomorrow.getDate();
    const weekdayFormatter = new Intl.DateTimeFormat('ja-JP', {
      timeZone: 'Asia/Tokyo',
      weekday: 'short',
    });
    const weekday = weekdayFormatter.format(tomorrow);
    return `${month}/${day}(${weekday})`;
  }

  const todayTitle = `きょう ${getTodayDateString()}`;
  const tomorrowTitle = `あした ${getTomorrowDateString()}`;
</script>

<main>
  <Header />
  <div class="content">
    <div class="left-column">
      <div class="calendar-today-tomorrow">
        <div class="calendar-today">
          <Calendar daysToShow={1} skipDays={0} title={todayTitle} />
        </div>
        <div class="calendar-tomorrow">
          <Calendar daysToShow={1} skipDays={1} title={tomorrowTitle} />
        </div>
      </div>
      <div class="calendar-upcoming">
        <Calendar daysToShow={5} skipDays={2} title="こんごのよてい" showDate={true} />
      </div>
    </div>
    <div class="right-column">
      <div class="right-top">
        <Weather />
      </div>
      <div class="right-bottom">
        <Tasks />
      </div>
    </div>
  </div>
</main>

<style>
  :global(body) {
    margin: 0;
    padding: 0;
    overflow: hidden;
  }

  :global(#app) {
    width: 100vw;
    height: 100vh;
    display: flex;
    flex-direction: column;
  }

  main {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    background: #ffffff;
  }

  .content {
    width: 100%;
    height: 95%;
    display: flex;
    gap: 8px;
    padding: 8px;
    box-sizing: border-box;
    background: #f0f4f8;
  }

  .left-column {
    width: 60%;
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .calendar-today-tomorrow {
    width: 100%;
    height: 60%;
    display: flex;
    gap: 8px;
  }

  .calendar-today {
    width: 60%;
    display: flex;
  }

  /* 今日エリアの左右padding調整 */
  :global(.calendar-today .calendar-widget) {
    padding: 8px 12px 24px 12px;
  }

  .calendar-tomorrow {
    width: 40%;
    display: flex;
  }

  /* 明日エリアの予定を20%小さく */
  :global(.calendar-tomorrow .calendar-widget) {
    padding: 8px 12px 24px 12px;
  }

  :global(.calendar-tomorrow .event) {
    font-size: 1.04rem;
    padding: 6px 8px;
  }

  :global(.calendar-tomorrow .event-time) {
    font-size: 0.6rem;
  }

  :global(.calendar-tomorrow .event-title) {
    font-size: 1.04rem;
  }

  .calendar-upcoming {
    width: 100%;
    height: 40%;
    display: flex;
  }

  /* 今後の予定エリアを3列表示に */
  :global(.calendar-upcoming .calendar-widget) {
    padding: 8px 12px 12px 12px;
  }

  :global(.calendar-upcoming .widget-title) {
    margin: 0 0 6px 0;
    font-size: 1.3rem;
  }

  :global(.calendar-upcoming .days-container) {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    grid-template-rows: repeat(6, auto);
    grid-auto-flow: column;
    gap: 6px;
    column-gap: 10px;
  }

  /* 日ごとのエリアを透明化してイベントを直接配置 */
  :global(.calendar-upcoming .day) {
    display: contents;
  }

  :global(.calendar-upcoming .events-container) {
    display: contents;
  }

  :global(.calendar-upcoming .all-day-section),
  :global(.calendar-upcoming .timed-events) {
    display: contents;
  }

  :global(.calendar-upcoming .event) {
    font-size: 0.85rem;
    padding: 4px 6px;
    gap: 6px;
    border-left-width: 4px;
  }

  :global(.calendar-upcoming .event-date) {
    font-size: 0.75rem;
    min-width: 30px;
  }

  :global(.calendar-upcoming .event-time) {
    font-size: 0.55rem;
    min-width: 32px;
  }

  :global(.calendar-upcoming .event-title) {
    font-size: 0.85rem;
  }

  .right-column {
    width: 40%;
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .right-top {
    width: 100%;
    height: 50%;
    display: flex;
  }

  .right-bottom {
    width: 100%;
    height: 50%;
    display: flex;
  }
</style>
