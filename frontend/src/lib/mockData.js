/**
 * フロント単体動作向けのモックデータ
 */

const MOCK_STATUS_OFFLINE = import.meta.env.VITE_MOCK_STATUS_OFFLINE === 'true';

function getTokyoDateParts(date) {
  const formatter = new Intl.DateTimeFormat('ja-JP', {
    timeZone: 'Asia/Tokyo',
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  });
  const parts = formatter.formatToParts(date);
  const year = parts.find((p) => p.type === 'year')?.value;
  const month = parts.find((p) => p.type === 'month')?.value;
  const day = parts.find((p) => p.type === 'day')?.value;
  return { year, month, day };
}

function formatTokyoYmd(date) {
  const { year, month, day } = getTokyoDateParts(date);
  return `${year}-${month}-${day}`;
}

function formatTokyoDateTime(date, hour, minute) {
  const ymd = formatTokyoYmd(date);
  const hh = String(hour).padStart(2, '0');
  const mm = String(minute).padStart(2, '0');
  return `${ymd}T${hh}:${mm}:00+09:00`;
}

function buildCalendarDays() {
  const base = new Date();
  const days = [];
  for (let i = 0; i < 7; i += 1) {
    const date = new Date(base);
    date.setDate(date.getDate() + i);
    const ymd = formatTokyoYmd(date);
    const day = {
      date: ymd,
      allDay: [],
      timed: [],
    };

    if (i === 0) {
      day.allDay.push({
        id: 'ad-1',
        title: '家族会議',
        start: ymd,
        end: ymd,
        color: '#60a5fa',
        calendar: 'Family',
      });
      day.timed.push({
        id: 't-1',
        title: '保育園お迎え',
        start: formatTokyoDateTime(date, 17, 30),
        end: formatTokyoDateTime(date, 18, 0),
        color: '#34d399',
        calendar: 'Family',
      });
    }

    if (i === 1) {
      day.timed.push({
        id: 't-2',
        title: '買い出し',
        start: formatTokyoDateTime(date, 10, 0),
        end: formatTokyoDateTime(date, 11, 0),
        color: '#fbbf24',
        calendar: 'Family',
      });
      day.timed.push({
        id: 't-3',
        title: 'オンライン面談',
        start: formatTokyoDateTime(date, 20, 0),
        end: formatTokyoDateTime(date, 20, 30),
        color: '#f87171',
        calendar: 'Family',
      });
    }

    if (i === 3) {
      day.allDay.push({
        id: 'ad-2',
        title: '通院',
        start: ymd,
        end: ymd,
        color: '#a78bfa',
        calendar: 'Family',
      });
    }

    if (i === 4) {
      day.timed.push({
        id: 't-4',
        title: '習い事',
        start: formatTokyoDateTime(date, 16, 0),
        end: formatTokyoDateTime(date, 17, 0),
        color: '#fb7185',
        calendar: 'Family',
      });
    }

    days.push(day);
  }
  return days;
}

function buildTasks() {
  const base = new Date();
  const dueToday = formatTokyoYmd(base);
  const dueTomorrow = formatTokyoYmd(new Date(base.getTime() + 24 * 60 * 60 * 1000));
  const dueLater = formatTokyoYmd(new Date(base.getTime() + 3 * 24 * 60 * 60 * 1000));
  const dueNext = formatTokyoYmd(new Date(base.getTime() + 5 * 24 * 60 * 60 * 1000));
  const dueNextWeek = formatTokyoYmd(new Date(base.getTime() + 7 * 24 * 60 * 60 * 1000));

  return [
    {
      id: 'task-1',
      title: 'ゴミ出し（燃えるゴミ）',
      dueDate: dueToday,
      priority: 'HIGH',
      completed: false,
    },
    {
      id: 'task-2',
      title: '夕飯の献立決める',
      dueDate: dueTomorrow,
      priority: 'MEDIUM',
      completed: false,
    },
    {
      id: 'task-3',
      title: '牛乳と卵を買う',
      dueDate: dueLater,
      priority: 'LOW',
      completed: false,
    },
    {
      id: 'task-4',
      title: '洗濯ネットを交換',
      dueDate: null,
      priority: 'LOW',
      completed: true,
    },
    {
      id: 'task-5',
      title: '授業参観の確認',
      dueDate: dueLater,
      priority: 'HIGH',
      completed: false,
    },
    {
      id: 'task-6',
      title: '水筒のパッキン洗浄',
      dueDate: dueTomorrow,
      priority: 'MEDIUM',
      completed: false,
    },
    {
      id: 'task-7',
      title: '薬の補充',
      dueDate: dueLater,
      priority: 'HIGH',
      completed: false,
    },
    {
      id: 'task-8',
      title: '給食セットの確認',
      dueDate: dueTomorrow,
      priority: 'MEDIUM',
      completed: false,
    },
    {
      id: 'task-9',
      title: 'プリント整理',
      dueDate: null,
      priority: 'LOW',
      completed: false,
    },
    {
      id: 'task-10',
      title: '日用品ストック確認',
      dueDate: dueNext,
      priority: 'LOW',
      completed: false,
    },
    {
      id: 'task-11',
      title: 'クリーニング受け取り',
      dueDate: dueToday,
      priority: 'HIGH',
      completed: false,
    },
    {
      id: 'task-12',
      title: '水道料金の支払い',
      dueDate: dueNextWeek,
      priority: 'MEDIUM',
      completed: false,
    },
    {
      id: 'task-13',
      title: '非常食の補充メモ',
      dueDate: null,
      priority: 'LOW',
      completed: true,
    },
  ];
}

function buildWeather() {
  const base = new Date();
  const weekly = [];
  const weeklyIcons = ['01d', '02d', '03d', '10d', '01d', '02d', '13d'];
  for (let i = 0; i < 7; i += 1) {
    const date = new Date(base);
    date.setDate(date.getDate() + i);
    weekly.push({
      date: formatTokyoYmd(date),
      icon: weeklyIcons[i % weeklyIcons.length],
      maxTemp: 10 + i,
      minTemp: 3 + Math.max(0, i - 2),
    });
  }

  return {
    location: '姫路市',
    current: {
      temperature: 12.3,
      condition: '曇',
      icon: '03d',
      humidity: 62,
      windSpeed: 2.4,
    },
    today: {
      maxTemp: 14.0,
      minTemp: 7.2,
      summary: 'くもり時々晴れ',
    },
    precipSlots: [
      { time: '06:00', probability: 10, icon: '01d' },
      { time: '09:00', probability: 20, icon: '02d' },
      { time: '12:00', probability: 30, icon: '03d' },
      { time: '15:00', probability: 40, icon: '10d' },
      { time: '18:00', probability: 25, icon: '03d' },
      { time: '21:00', probability: 15, icon: '02d' },
      { time: '00:00', probability: 10, icon: '01n' },
      { time: '03:00', probability: 5, icon: '01n' },
    ],
    weekly,
    alerts: [
      {
        severity: '特別警報',
        title: '大雨特別警報',
      },
    ],
  };
}

function buildStatus() {
  const now = new Date().toISOString();
  return {
    ok: true,
    now,
    errors: [],
    lastUpdated: {
      weather: now,
      calendar: now,
      tasks: now,
    },
  };
}

export function getMockResponse(endpoint) {
  switch (endpoint) {
    case '/api/status':
      if (MOCK_STATUS_OFFLINE) {
        throw new Error('Mock status offline');
      }
      return buildStatus();
    case '/api/calendar':
      return { days: buildCalendarDays() };
    case '/api/tasks':
      return { items: buildTasks() };
    case '/api/weather':
      return buildWeather();
    default:
      return {};
  }
}
