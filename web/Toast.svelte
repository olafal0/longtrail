<script>
  import { createEventDispatcher } from "svelte";

  const dispatch = createEventDispatcher();

  let messages = [];

  export function addMsg(text, alertType) {
    alertType = alertType || "alert-primary";
    // Create a message object that can remove itself
    const msg = {
      text,
      type: alertType,
      remove: () => {
        const idx = messages.indexOf(msg);
        if (idx >= 0) {
          removeMsg(idx);
        }
      }
    };
    // Push new message, preserving reactivity
    messages = [...messages, msg];
    // Automatically remove after 6 seconds (0.5s fadein/fadeout + 5s visibility)
    setTimeout(msg.remove, 6000);
    return msg;
  }

  export function addError(text) {
    return addMsg(text, "alert-danger");
  }

  function removeMsg(index) {
    // equivalent of splice, but using assignment to make sure Svelte reacts to
    // the change correctly
    messages = [
      ...messages.slice(0, index),
      ...messages.slice(index + 1, messages.length)
    ];
  }
</script>

<style>
  @keyframes fadein {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }

  @keyframes fadeout {
    from {
      opacity: 1;
    }
    to {
      opacity: 0;
    }
  }

  .popup {
    animation: fadein 0.5s, fadeout 0.5s 4s forwards;
    padding: 10px;
    margin-bottom: 10px;
    border-radius: 5px;
  }

  .popup-container {
    position: fixed;
    left: 50%;
    top: 5px;
    z-index: 1;
  }
</style>

{#if messages}
  <div class="popup-container">
    {#each messages as msg}
      <div class="popup {msg.type}" on:click={msg.remove}>{msg.text}</div>
    {/each}
  </div>
{/if}
