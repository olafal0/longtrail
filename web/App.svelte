<script>
  import Auth from "./Auth";
  import CalendarView from "./CalendarView.svelte";
  import Login from "./Login.svelte";

  let user = null;

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

<div>
  <nav class="navbar bg-primary">
    <button class="btn btn-clear text-light">Longtrail</button>
    {#if user}
      <div class="right">
        <span class="text-light">{user.username}</span>
        <button class="btn bg-light" on:click={logout}>Log Out</button>
      </div>
    {/if}
  </nav>
  {#if user}
    <CalendarView />
  {:else}
    <Login on:signedIn={signedIn} />
  {/if}
</div>
