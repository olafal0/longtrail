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

  let editingEvent = null;
  let editingEventTitle = "";

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

  async function deleteEvent() {
    try {
      await api.deletePractice(editingEvent.id);
      editingEvent.remove();
    } catch (e) {
      toasts.addError(e);
    } finally {
      editingEvent = null;
    }
  }

  function eventClick({ detail: eventElement }) {
    // Set the clicked event as the currently editing event
    const event = eventElement.event;
    editingEventTitle = event.title;
    editingEvent = event;
  }

  async function updateEditedEvent() {
    // Updating the title requires editing the title in FullCalendar as well
    // as saving the title to additionalData
    editingEvent.setProp("title", editingEventTitle);
    editingEvent.setExtendedProp("additionalData", {
      ...editingEvent.additionalData,
      title: editingEventTitle
    });
    await updateEvent({ detail: { event: editingEvent } });
    editingEvent = null;
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

  function shortTime(date) {
    const dayFormatOptions = {
      weekday: "short",
      hour: "numeric",
      minute: "numeric",
      timeZoneName: "short"
    };
    return new Intl.DateTimeFormat("en-US", dayFormatOptions).format(date);
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

{#if editingEvent}
  <Modal
    on:close={() => {
      editingEvent = null;
    }}>
    <div class="card">
      <div class="card-body">
        <h5 class="card-title">Edit Event</h5>
        <p>
          <i>
            From {shortTime(editingEvent.start)} to {shortTime(editingEvent.end)}
          </i>
        </p>
        <form on:submit|preventDefault={updateEditedEvent}>
          <div>
            <label>
              Title:
              <p>
                <input type="text" bind:value={editingEventTitle} />
              </p>
            </label>
          </div>
          <input
            type="button"
            class="btn btn-secondary"
            on:click={() => {
              editingEvent = null;
            }}
            value="Cancel" />
          <input
            type="button"
            class="btn btn-danger"
            on:click={deleteEvent}
            value="Delete" />
          <input type="submit" class="btn btn-primary" value="Update" />
        </form>
      </div>
    </div>
  </Modal>
{/if}

<Toast bind:this={toasts} />
