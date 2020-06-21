<script>
  import Auth from "./Auth";
  import CalendarView from "./CalendarView.svelte";
  import Login from "./Login.svelte";
  import Modal from "./Modal.svelte";

  let user = null;
  let navbarActions = [
    {
      text: "Log Out",
      clicked: logout
    }
  ];
  let calendarActions = [];
  let showAbout = false;

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
    <div class="text-light navbar-brand">Long Trail</div>
    <div
      class="btn btn-clear text-light mr-auto"
      on:click={() => {
        showAbout = true;
      }}>
      About
    </div>
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

{#if showAbout}
  <Modal
    on:close={() => {
      showAbout = false;
    }}>
    <div class="card">
      <div class="card-body">
        <h5 class="card-title">Long Trail: About</h5>
        <p>
          Long Trail was designed to help
          <a href="https://gmcmf.org/content/about-green-mountain">
            Green Mountain Chamber Music Festival
          </a>
          students stay on top of practicing. The name "Long Trail" is a
          reference both to the Vermont portion of the Appalachian Trail, and
          also to the observation of Ivan Galamian that “the road to mastery [of
          our instrument] is a long and arduous one.”
        </p>
        <p>
          Long Trail was built by Daniel Lawrence. This website and its
          infrastructure are entirely open source and MIT-licensed. You can view
          the source code here:
          <br />
          <a href="https://github.com/olafal0/longtrail">
            https://github.com/olafal0/longtrail
          </a>
        </p>
        <p>
          Thank you to Kevin and Barbara Lawrence, without whom I would not be
          where I am today. Much love, and happy father's day ♥
        </p>
        <button
          class="btn btn-primary mx-auto"
          on:click={() => {
            showAbout = false;
          }}>
          Close
        </button>
      </div>
    </div>
  </Modal>
{/if}
