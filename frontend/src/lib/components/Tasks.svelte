<script>
  import { onMount } from 'svelte';
  import { getTasks } from '../api.js';

  let tasksData = null;
  let error = null;

  const MAX_DISPLAY_TASKS = 12;

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
    return () => clearInterval(interval);
  });
</script>

<div class="tasks-widget">
  {#if error}
    <div class="error">エラー: {error}</div>
  {:else if tasksData && tasksData.items}
    <h2 class="widget-title">タスク一覧</h2>
    {#if tasksData.items.length === 0}
      <div class="no-tasks">タスクはありません</div>
    {:else}
      <div class="tasks-list">
        {#each tasksData.items.slice(0, MAX_DISPLAY_TASKS) as task, idx}
          <div class="task-item" style={getTaskStyle(task)}>
            <div class="task-content">
              <div class="task-title">{task.title}</div>
              <div class="task-meta">
                <span class="task-due">{formatDueDate(task.dueDate)}</span>
                {#if task.priority}
                  <span class="task-priority">{getPriorityLabel(task.priority)}</span>
                {/if}
              </div>
            </div>
            {#if task.completed}
              <div class="task-badge badge-completed">完了</div>
            {/if}
          </div>
        {/each}
      </div>
      {#if tasksData.items.length > MAX_DISPLAY_TASKS}
        <div class="more-tasks">
          他 {tasksData.items.length - MAX_DISPLAY_TASKS} 件
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
    font-size: 1.2rem;
    padding: 40px;
  }

  .error {
    color: #ef4444;
  }

  .widget-title {
    margin: 0 0 16px 0;
    font-size: 1.6rem;
    font-weight: bold;
    color: #1f2937;
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
    gap: 10px;
    flex: 1;
  }

  .task-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 16px;
    border-left: 4px solid #d1d5db;
    background: #f9fafb;
    border-radius: 4px;
    transition: all 0.2s ease;
  }

  .task-item:hover {
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }

  .task-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    justify-content: center;
    gap: 6px;
  }

  .task-title {
    font-size: 1.1rem;
    font-weight: 500;
    color: #1f2937;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .task-meta {
    display: flex;
    gap: 12px;
    font-size: 0.9rem;
    color: #6b7280;
  }

  .task-due {
    font-weight: 500;
  }

  .task-priority {
    font-size: 1rem;
  }

  .task-badge {
    padding: 4px 12px;
    border-radius: 12px;
    font-size: 0.8rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.5px;
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
