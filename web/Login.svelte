<script>
  import Auth from "./Auth";
  import { onMount, createEventDispatcher } from "svelte";
  import Form from "./Form.svelte";

  const dispatch = createEventDispatcher();

  const FlowState = Object.freeze({
    loading: 1,
    initial: 2,
    signUpConfirmation: 3
  });

  const initialFormFields = {
    email: {
      type: "email",
      label: "Email Address"
    },
    password: {
      type: "password",
      label: "Password"
    }
  };

  const signUpConfirmationFormFields = {
    code: {
      label: "Confirmation Code",
      type: "text"
    }
  };

  let mainForm = null;
  let _username = "";
  let flowState = FlowState.loading;
  let infoMessage = "";
  let errorMessage = "";

  onMount(() => {
    Auth.currentAuthenticatedUser()
      .then(signedIn)
      .catch(err => {
        if (err === "not authenticated") {
          flowState = FlowState.initial;
        } else {
          displayError(err);
        }
      });
  });

  function msgToString(msg) {
    if (typeof msg === "string") {
      return msg;
    }
    if (msg.message && typeof msg.message === "string") {
      return msg.message;
    }
    return JSON.stringify(msg);
  }

  function displayError(msg) {
    errorMessage = msgToString(msg);
    infoMessage = "";
  }

  function displayInfo(msg) {
    infoMessage = msgToString(msg);
    errorMessage = "";
  }

  function signedIn() {
    dispatch("signedIn");
  }

  function signUp({ email, password }) {
    displayInfo("Registering...");
    _username = email;
    Auth.signUp({
      username: email,
      password,
      attributes: {
        email
      }
    })
      .then(({ user, userConfirmed, userSub }) => {
        if (!userConfirmed) {
          // User is not confirmed; we need a confirmation code
          displayInfo(
            "Please confirm by entering the confirmation code sent to your email address"
          );
          flowState = FlowState.signUpConfirmation;
          return;
        }
        signedIn();
      })
      .catch(displayError);
  }

  function resendCode() {
    displayInfo(`Resending confirmation code for ${_username}`);
    Auth.resendSignUp(_username).catch(displayError);
  }

  function confirmSignUp({ detail: data }) {
    displayInfo("Confirming user...");
    Auth.confirmSignUp(_username, data.code)
      .then(result => {
        if (result === "SUCCESS") {
          // Set back to initial, now that user needs to sign in for the first time
          displayInfo("User confirmed! Please sign in");
          flowState = FlowState.initial;
          return;
        }
      })
      .catch(displayError);
  }

  function signIn({ detail: data }) {
    displayInfo("Signing in...");
    _username = data.email;
    Auth.signIn(data.email, data.password)
      .then(signedIn)
      .catch(err => {
        if (err.code === "UserNotConfirmedException") {
          displayInfo(
            "Please confirm by entering the confirmation code sent to your email address"
          );
          flowState = FlowState.signUpConfirmation;
        } else {
          displayError(err);
        }
      });
  }
</script>

<div class="container">
  <div class="card">
    <div class="card-body">
      {#if errorMessage}
        <div class="alert alert-danger">{errorMessage}</div>
      {:else if infoMessage}
        <div class="alert alert-primary">{infoMessage}</div>
      {/if}
      <h5 class="card-title">Log in</h5>
      {#if flowState === FlowState.initial}
        <Form
          bind:this={mainForm}
          formFields={initialFormFields}
          on:submit={signIn}>
          <button
            class="btn btn-secondary"
            on:click|preventDefault={() => signUp(mainForm.getFormData())}>
            Register
          </button>
          <button type="submit" class="btn btn-primary">Log In</button>
        </Form>
      {:else if flowState === FlowState.signUpConfirmation}
        <Form
          bind:this={mainForm}
          formFields={signUpConfirmationFormFields}
          on:submit={confirmSignUp}>
          <!-- TODO: Pressing return in this form resends the confirmation code,
          instead of calling confirmSignUp -->
          <button
            class="btn btn-secondary"
            on:click|preventDefault={() => resendCode(mainForm.getFormData())}>
            Resend Confirmation Code
          </button>
          <button type="submit" class="btn btn-primary">Confirm</button>
        </Form>
      {/if}
    </div>
  </div>
</div>
