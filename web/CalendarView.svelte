<script>
  import Auth from "./Auth";
  import FullCalendar from "./FullCalendar.svelte";
  import { onMount, tick } from "svelte";
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

  function loadEvents(info) {
    return api.getPractices(info.startStr, info.endStr);
  }

  async function dateClick({ detail: clicked }) {
    const start = clicked.date;
    const end = addMinutes(start, appConfig.defaultDuration);
    const event = {
      title: "Practice Session",
      start,
      end,
      startEditable: true,
      durationEditable: true
    };

    // Add an un-editable fake event first, to avoid a visual delay
    const addedEvent = fc.addEvent({ start, end, title: "Saving..." });

    // Create event through API to get event ID
    const eventId = await api.createPractice(event);

    // Remove fc event and add real one, with correct ID
    addedEvent.remove();
    fc.addEvent(event);
  }

  function eventClick({ detail: event }) {
    event.event.remove();
  }
</script>

<FullCalendar
  config={fcConfig}
  bind:this={calendarComponent}
  on:dateClick={dateClick}
  on:eventClick={eventClick} />
