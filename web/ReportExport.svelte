<script>
  import Modal from "./Modal";
  import { onMount, createEventDispatcher } from "svelte";
  import { addMinutes, addDays } from "./vanillaDate";
  import ClipboardJS from "clipboard";

  export let reportEvents;
  let reportTextArea;
  let copyButton;
  let errorText = "";

  const dispatch = createEventDispatcher();

  onMount(() => {
    const cb = new ClipboardJS(copyButton);
  });
</script>

<Modal on:close>
  <div class="card">
    <div class="card-body">
      <h5 class="card-title">Hours practiced (last seven days)</h5>
      <p id="reportText" bind:this={reportTextArea}>
        {#each reportEvents as day}
          <div>{day.dateStr}: {day.durationStr}</div>
        {/each}
      </p>
      <button class="btn btn-secondary" on:click={() => dispatch('close')}>
        Cancel
      </button>
      <button
        bind:this={copyButton}
        class="btn btn-primary"
        data-clipboard-target="#reportText">
        Copy to Clipboard
      </button>
    </div>
  </div>
</Modal>
