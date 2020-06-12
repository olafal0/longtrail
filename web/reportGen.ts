interface AdditionalData {
  title: string;
}

interface CalEvent {
  additionalData: AdditionalData;
  end: string;
  id: string;
  start: string;
  userId: string;
}

interface DayTotal {
  date: number;
  dateStr: string;
  durationStr: string;
}

const dayFormatOptions = {
  weekday: 'short',
  month: 'long',
  day: 'numeric',
};

const msToMinutes = 1 / (1000 * 60);

export default function generateTextReport(events: CalEvent[]): DayTotal[] {
  events = events.sort((a, b) => a.start.localeCompare(b.start));
  // Iterate through events
  // For each event, add its duration to the total time practiced that day
  const minutesByDay: object = {};
  events.forEach((e) => {
    const start = new Date(e.start);
    const end = new Date(e.end);
    const startDay = new Date(start);
    startDay.setHours(0);
    startDay.setMinutes(0);
    startDay.setSeconds(0);
    startDay.setMilliseconds(0);

    const durationMinutes = Math.round(
      (end.getTime() - start.getTime()) * msToMinutes
    );

    minutesByDay[startDay.toISOString()] =
      (minutesByDay[startDay.toISOString()] || 0) + durationMinutes;
  });

  // Convert hours by day into list of {date, dateStr, durationStr}
  const dayTotals = [];
  Object.getOwnPropertyNames(minutesByDay).forEach((dayStr) => {
    const date = new Date(dayStr);
    const dateStr = new Intl.DateTimeFormat('en-US', dayFormatOptions).format(
      date
    );

    // Create hours and minutes strings, e.g. "02"
    const durationMinutes = minutesByDay[dayStr] || 0;
    const minutes = durationMinutes % 60;
    const hours = (durationMinutes - minutes) / 60;
    const hoursStr = hours.toFixed(0).padStart(2, '0');
    const minutesStr = minutes.toFixed(0).padStart(2, '0');
    const durationStr = `${hoursStr}:${minutesStr}`;

    dayTotals.push({
      date: date.getTime(),
      dateStr,
      durationStr,
    });
  });

  dayTotals.sort((a, b) => a.date - b.date);

  return dayTotals;
}
