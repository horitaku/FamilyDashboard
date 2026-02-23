<script>
  import { onMount } from 'svelte';
  import { getTasks } from '../api.js';

  let tasksData = null;
  let error = null;
  let tasksWidgetEl = null;
  let displayCount = 12;

  const ESTIMATED_ROW_HEIGHT = 46;
  const WIDGET_VERTICAL_PADDING = 16;

  /**
   * タスクデータを取得
   */
  async function loadTasksData() {
    try {
      tasksData = await getTasks();
      error = null;
    } catch (err) {
      console.error('タスクデータ取得エラー:', err);
      error = err.message;
    }
  }

  /**
   * 期限フォーマット（MM月DD日）
   */
  function formatDueDate(dueDate) {
    if (!dueDate) return '期限なし';
    const date = new Date(dueDate + 'T00:00');
    const formatter = new Intl.DateTimeFormat('ja-JP', {
      month: 'numeric',
      day: 'numeric',
    });
    return formatter.format(date);
  }

  /**
   * 期限までの日数を計算
   */
  function getDaysUntilDue(dueDate) {
    if (!dueDate) return Infinity;
    const today = new Date();
    today.setHours(0, 0, 0, 0);
    const dueDateObj = new Date(dueDate + 'T00:00');
    const diffMs = dueDateObj.getTime() - today.getTime();
    return Math.ceil(diffMs / (1000 * 60 * 60 * 24));
  }

  /**
   * タスク行のスタイルを生成（優先度により）
   */
  function getTaskStyle(task) {
    let borderColor = '#d1d5db'; // グレー
    let backgroundColor = '#f9fafb';

    // 優先度により色を変更
    if (task.priority === 'HIGH') {
      borderColor = '#ef4444';
      backgroundColor = '#fee2e2';
    } else if (task.priority === 'MEDIUM') {
      borderColor = '#f59e0b';
      backgroundColor = '#fef3c7';
    } else if (task.priority === 'LOW') {
      borderColor = '#10b981';
      backgroundColor = '#ecfdf5';
    }

    // 期限が過ぎている場合は赤く表示
    const daysUntil = getDaysUntilDue(task.dueDate);
    if (daysUntil < 0) {
      borderColor = '#dc2626';
      backgroundColor = '#fecaca';
    }

    return `border-left-color: ${borderColor}; background-color: ${backgroundColor}`;
  }

  /**
   * 優先度を日本語で表示
   */
  function getPriorityLabel(priority) {
    switch (priority) {
      case 'HIGH':
        return '⭐⭐⭐';
      case 'MEDIUM':
        return '⭐⭐';
      case 'LOW':
        return '⭐';
      default:
        return '';
    }
  }

  onMount(() => {
    loadTasksData();
    // 5分ごとにリロード
    const interval = setInterval(loadTasksData, 300000);
    let resizeObserver = null;

    if (tasksWidgetEl && 'ResizeObserver' in window) {
      resizeObserver = new ResizeObserver((entries) => {
        const entry = entries[0];
        if (!entry) return;
        const availableHeight = Math.max(
          0,
          entry.contentRect.height - WIDGET_VERTICAL_PADDING
        );
        const nextCount = Math.max(
          1,
          Math.floor(availableHeight / ESTIMATED_ROW_HEIGHT)
        );
        displayCount = nextCount;
      });
      resizeObserver.observe(tasksWidgetEl);
    }

    return () => {
      clearInterval(interval);
      if (resizeObserver) resizeObserver.disconnect();
    };
  });
</script>

<div class="tasks-widget" bind:this={tasksWidgetEl}>
  {#if error}
    <div class="error">エラー: {error}</div>
  {:else if tasksData && tasksData.items}
    {#if tasksData.items.length === 0}
      <div class="no-tasks">タスクはありません</div>
    {:else}
      <div class="tasks-list">
        {#each tasksData.items.slice(0, displayCount) as task, idx}
          <div class="task-item" style={getTaskStyle(task)}>
            <div class="task-line">
              <span class="task-title">{task.title}</span>
              <span class="task-due">{formatDueDate(task.dueDate)}</span>
              {#if task.priority}
                <span class="task-priority">{getPriorityLabel(task.priority)}</span>
              {/if}
              {#if task.completed}
                <span class="task-badge badge-completed">完了</span>
              {/if}
            </div>
          </div>
        {/each}
      </div>
      {#if tasksData.items.length > displayCount}
        <div class="more-tasks">
          他 {tasksData.items.length - displayCount} 件
        </div>
      {/if}
    {/if}
  {:else}
    <div class="loading">読み込み中...</div>
  {/if}
</div>

<style>
  .tasks-widget {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    padding: 8px;
    box-sizing: border-box;
    background: linear-gradient(135deg, #f0f3f8 0%, #e0e6f0 100%);
    border-radius: 8px;
    overflow: hidden;
  }

  .error,
  .loading {
    text-align: center;
    color: #666;
    font-size: 1.2rem;
    padding: 40px;
  }

  .error {
    color: #ef4444;
  }

  .no-tasks {
    text-align: center;
    color: #9ca3af;
    font-size: 1.1rem;
    padding: 40px 0;
  }

  .tasks-list {
    display: flex;
    flex-direction: column;
    gap: 6px;
    flex: 1;
  }

  .task-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 6px 10px;
    border-left: 8px solid #d1d5db;
    background: #f9fafb;
    border-radius: 4px;
    transition: all 0.2s ease;
  }

  .task-item:hover {
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }

  .task-line {
    display: flex;
    align-items: center;
    gap: 10px;
    width: 100%;
    min-width: 0;
  }

  .task-title {
    font-size: 1rem;
    font-weight: 500;
    color: #1f2937;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    flex: 1;
  }

  .task-due {
    font-weight: 500;
    font-size: 0.85rem;
    color: #6b7280;
    flex-shrink: 0;
  }

  .task-priority {
    font-size: 0.9rem;
    flex-shrink: 0;
  }

  .task-badge {
    padding: 2px 8px;
    border-radius: 12px;
    font-size: 0.7rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    flex-shrink: 0;
  }

  .badge-completed {
    background: #d1fae5;
    color: #065f46;
  }

  .more-tasks {
    margin-top: 12px;
    padding: 12px;
    text-align: center;
    background: rgba(156, 163, 175, 0.1);
    border-radius: 4px;
    color: #6b7280;
    font-size: 0.95rem;
  }

  /* スクロールバーのスタイル */
  ::-webkit-scrollbar {
    width: 8px;
  }

  ::-webkit-scrollbar-track {
    background: transparent;
  }

  ::-webkit-scrollbar-thumb {
    background: #d1d5db;
    border-radius: 4px;
  }

  ::-webkit-scrollbar-thumb:hover {
    background: #9ca3af;
  }
</style>
