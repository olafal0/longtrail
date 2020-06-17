<script>
  import Auth from "./Auth";
  import CalendarView from "./CalendarView.svelte";
  import Login from "./Login.svelte";

  let user = null;
  let navbarActions = [
    {
      text: "Log Out",
      clicked: logout
    }
  ];
  let calendarActions = [];

  function signedIn() {
    Auth.currentAuthenticatedUser()
      .then(userData => {
        user = userData;
      })
      .catch(console.error);
  }

  function logout() {
    Auth.signOut().then(() => {
      user = null;
    });
  }
</script>

<style>
  .nav-item {
    margin: 0 5px;
  }

  .navbar-brand {
    font-size: 1.75rem;
  }
</style>

<div class="app-container">
  <nav class="navbar bg-primary">
    <div class="btn btn-clear text-light navbar-brand">Long Trail</div>
    {#if user}
      <div class="right">
        <span class="text-light">{user.username}</span>
        {#each [...calendarActions, ...navbarActions] as action}
          <button class="nav-item btn bg-light" on:click={action.clicked}>
            {action.text}
          </button>
        {/each}
      </div>
    {/if}
  </nav>
  {#if user}
    <CalendarView bind:navbarActions={calendarActions} />
  {:else}
    <Login on:signedIn={signedIn} />
  {/if}
</div>
