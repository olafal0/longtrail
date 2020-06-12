<script>
  import FullCalendar from "./FullCalendar";
  import Modal from "./Modal";
  import Toast from "./Toast";
  import ReportExport from "./ReportExport";
  import { onMount } from "svelte";
  import api from "./api";
  import { addMinutes, addDays } from "./vanillaDate";
  import generateTextReport from "./reportGen";

  const fcConfig = {
    height: "auto",
    minTime: "06:00:00",
    slotDuration: "00:15:00",
    allDaySlot: false,
    header: {
      left: "prev,next today",
      center: "title",
      right: "timeGridWeek,timeGridDay"
    },
    events: loadEvents
  };

  const appConfig = {
    // Default duration of events in minutes
    defaultDuration: 60
  };

  export let navbarActions;
  let calendarComponent;
  let fc;
  let toasts;
  let showExportModal = false;
  let reportEvents;

  const exportLogNavbarAction = {
    text: "Export Log",
    created: null,
    clicked: generateLog
  };

  onMount(() => {
    fc = calendarComponent.getCalendar();
    navbarActions = [exportLogNavbarAction];
  });

  async function loadEvents({ startStr, endStr }) {
    // Remove existing events. This prevents seeing duplicates when an event is
    // recently created, but also gets loaded from the backend. If fc isn't set
    // yet, we don't need to worry about duplicates.
    if (fc) {
      fc.getEvents().forEach(e => e.remove());
    }
    const events = await api.getPractices(startStr, endStr);
    const eventObjects = events.map(e => ({
      startEditable: true,
      durationEditable: true,
      ...e,
      ...(e.additionalData || {})
    }));
    return eventObjects;
  }

  async function dateClick({ detail: clicked }) {
    const start = clicked.date;
    const end = addMinutes(start, appConfig.defaultDuration);
    const event = {
      title: "Practice Session",
      start,
      end,
      startEditable: true,
      durationEditable: true,
      additionalData: {
        title: "Practice Session"
      }
    };

    // Add an un-editable fake event first, to avoid a visual delay
    const addedEvent = fc.addEvent({ start, end, title: "Saving..." });

    // Create event through API to get event ID
    try {
      const eventId = await api.createPractice(event);
      event.id = eventId;
      // Add real event with correct ID
      fc.addEvent(event);
    } catch (e) {
      toasts.addError(e);
    } finally {
      // Remove fake placeholder event
      addedEvent.remove();
    }
  }

  async function eventClick({ detail: eventElement }) {
    // https://fullcalendar.io/docs/eventRender
    // hook for modifying event DOM
    const event = eventElement.event;
    try {
      await api.deletePractice(event.id);
      event.remove();
    } catch (e) {
      toasts.addError(e);
    }
  }

  async function updateEvent({ detail: eventDrop }) {
    const event = {
      ...eventDrop.event.extendedProps,
      id: eventDrop.event.id,
      start: eventDrop.event.start,
      end: eventDrop.event.end
    };
    try {
      await api.setPractice(event);
    } catch (e) {
      toasts.addError(e);
      eventDrop.revert();
    }
  }

  async function generateLog() {
    const now = new Date();
    const startStr = addDays(now, -7).toISOString();
    const endStr = now.toISOString();
    const loadingNotification = toasts.addMsg("Loading...");
    try {
      const events = await api.getPractices(startStr, endStr);
      reportEvents = generateTextReport(events);
      showExportModal = true;
    } catch (e) {
      toasts.addError(e);
    } finally {
      loadingNotification.remove();
    }
  }
</script>

<FullCalendar
  config={fcConfig}
  bind:this={calendarComponent}
  on:dateClick={dateClick}
  on:eventClick={eventClick}
  on:eventDrop={updateEvent}
  on:eventResize={updateEvent} />

{#if showExportModal}
  <ReportExport {reportEvents} on:close={() => (showExportModal = false)} />
{/if}

<Toast bind:this={toasts} />
