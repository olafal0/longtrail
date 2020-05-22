<script>
  import Auth from "./Auth";
  import FullCalendar from "./FullCalendar.svelte";
  import Login from "./Login.svelte";
  import { onMount, tick } from "svelte";

  const fcConfig = {
    height: "auto",
    minTime: "07:00:00",
    allDaySlot: false
  };

  let user = null;
  let loggedIn = false;
  let calendarComponent;
  let fc;

  function signedIn() {
    Auth.currentAuthenticatedUser()
      .then(userData => {
        user = userData;
        loggedIn = true;
        tick().then(() => {
          fc = calendarComponent.getCalendar();
        });
      })
      .catch(console.error);
  }

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

<div>
  <nav class="navbar bg-primary">
    <button class="btn btn-clear text-light">Longtrail</button>
    {#if loggedIn}
      <div class="right">
        <span class="text-light">{user.username}</span>
        <button class="btn bg-light">Log Out</button>
      </div>
    {/if}
  </nav>
  {#if loggedIn}
    <FullCalendar
      config={fcConfig}
      bind:this={calendarComponent}
      on:dateClick={dateClick}
      on:eventClick={eventClick} />
  {:else}
    <Login on:signedIn={signedIn} />
  {/if}
</div>
