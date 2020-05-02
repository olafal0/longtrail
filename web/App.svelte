<script>
  import FullCalendar from "./FullCalendar.svelte";
  import { onMount } from "svelte";

  const fcConfig = {
    height: "auto",
    minTime: "07:00:00",
    allDaySlot: false
  };

  let calendarComponent;
  let fc;

  onMount(() => {
    fc = calendarComponent.getCalendar();
  });

  function dateClick({ detail: event }) {
    event = {
      ...event,
      editable: true,
      startEditable: true,
      durationEditable: true
    };
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
