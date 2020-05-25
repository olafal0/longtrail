<script>
  import { Calendar } from "@fullcalendar/core";
  import timeGridPlugin from "@fullcalendar/timegrid";
  import interaction from "@fullcalendar/interaction";
  import bootstrapPlugin from "@fullcalendar/bootstrap";
  import { onMount, onDestroy, createEventDispatcher } from "svelte";

  import "@fullcalendar/core/main.css";
  import "@fullcalendar/timegrid/main.css";

  const dispatch = createEventDispatcher();

  let domCalendar;
  let calendar;
  export let config;

  onMount(() => {
    calendar = new Calendar(domCalendar, {
      plugins: [timeGridPlugin, interaction, bootstrapPlugin],
      themeSystem: "bootstrap",
      dateClick: event => dispatch("dateClick", event),
      eventClick: event => dispatch("eventClick", event),
      eventDragStart: event => dispatch("eventDragStart", event),
      eventDragStop: event => dispatch("eventDragStop", event),
      eventDrop: event => dispatch("eventDrop", event),
      eventLeave: event => dispatch("eventLeave", event),
      eventMouseEnter: event => dispatch("eventMouseEnter", event),
      eventMouseLeave: event => dispatch("eventMouseLeave", event),
      eventResize: event => dispatch("eventResize", event),
      eventResizeStart: event => dispatch("eventResizeStart", event),
      eventResizeStop: event => dispatch("eventResizeStop", event),
      ...config
    });
    calendar.render();
  });

  onDestroy(() => {
    calendar.destroy();
  });

  export function getCalendar() {
    return calendar;
  }
</script>

<div bind:this={domCalendar} />
