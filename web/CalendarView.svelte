<script>
  import FullCalendar from "./FullCalendar.svelte";
  import { onMount } from "svelte";
  import api from "./api";
  import { addMinutes } from "./vanillaDate";

  const fcConfig = {
    height: "auto",
    minTime: "07:00:00",
    allDaySlot: false,
    events: loadEvents
  };

  const appConfig = {
    // Default duration of events in minutes
    defaultDuration: 60
  };

  let calendarComponent;
  let fc;

  onMount(() => {
    fc = calendarComponent.getCalendar();
  });

  async function loadEvents(info) {
    const events = await api.getPractices(info.startStr, info.endStr);
    const eventObjects = events.map(e => ({
      startEditable: true,
      durationEditable: true,
      ...e,
      ...(e.additionalData || {})
    }));
    console.log(eventObjects);
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
    const eventId = await api.createPractice(event);
    event.id = eventId;

    // Remove fc event and add real one, with correct ID
    addedEvent.remove();
    fc.addEvent(event);
  }

  async function eventClick({ detail: eventElement }) {
    const event = eventElement.event;
    await api.deletePractice(event.id);
    event.remove();
  }

  async function updateEvent({ detail: eventDrop }) {
    console.log(eventDrop.event.extendedProps);
    const event = {
      ...eventDrop.event.extendedProps,
      id: eventDrop.event.id,
      start: eventDrop.event.start,
      end: eventDrop.event.end
    };
    try {
      await api.setPractice(event);
    } catch (err) {
      console.error(err);
      eventDrop.revert();
    }
  }
</script>

<FullCalendar
  config={fcConfig}
  bind:this={calendarComponent}
  on:dateClick={dateClick}
  on:eventClick={eventClick}
  on:eventDrop={updateEvent} />
