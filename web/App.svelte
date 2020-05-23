<script>
  import Auth from "./Auth";
  import FullCalendar from "./FullCalendar.svelte";
  import Login from "./Login.svelte";
  import { onMount, tick } from "svelte";
  import api from "./api";

  const fcConfig = {
    height: "auto",
    minTime: "07:00:00",
    allDaySlot: false
  };

  let user = null;
  let calendarComponent;
  let fc;

  function signedIn() {
    Auth.currentSession().then(console.log);
    Auth.currentAuthenticatedUser()
      .then(userData => {
        user = userData;
        tick().then(() => {
          fc = calendarComponent.getCalendar();
        });
      })
      .catch(console.error);
  }

  function logout() {
    Auth.signOut().then(() => {
      user = null;
    });
  }

  function dateClick({ detail: event }) {
    event = {
      ...event,
      editable: true,
      startEditable: true,
      durationEditable: true
    };
    fc.addEvent(event);
    api.createPractice(event).catch(console.error);
  }

  function eventClick({ detail: event }) {
    event.event.remove();
  }
</script>

<div>
  <nav class="navbar bg-primary">
    <button class="btn btn-clear text-light">Longtrail</button>
    {#if user}
      <div class="right">
        <span class="text-light">{user.username}</span>
        <button class="btn bg-light" on:click={logout}>Log Out</button>
        <button
          class="btn bg-light"
          on:click={() => {
            api.echo();
          }}>
          Echo Test
        </button>
      </div>
    {/if}
  </nav>
  {#if user}
    <FullCalendar
      config={fcConfig}
      bind:this={calendarComponent}
      on:dateClick={dateClick}
      on:eventClick={eventClick} />
  {:else}
    <Login on:signedIn={signedIn} />
  {/if}
</div>
