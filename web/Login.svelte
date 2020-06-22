<script>
  import Auth from "./Auth";
  import { onMount, createEventDispatcher } from "svelte";

  const dispatch = createEventDispatcher();

  const FlowState = Object.freeze({
    loading: 1,
    initial: 2,
    signUpConfirmation: 3,
    forgotPassword: 4
  });

  let formFields = {
    email: "",
    password: "",
    newPassword: "",
    confirmationCode: ""
  };

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

  function clearMessages() {
    infoMessage = "";
    errorMessage = "";
  }

  function displayError(msg) {
    clearMessages();
    errorMessage = msgToString(msg);
  }

  function displayInfo(msg) {
    clearMessages();
    infoMessage = msgToString(msg);
  }

  function signedIn() {
    dispatch("signedIn");
  }

  function signUp() {
    displayInfo("Registering...");
    // TODO: Special case handling for common errors or problems (eg "@gmail.con")
    formFields.email = formFields.email;
    Auth.signUp({
      username: formFields.email,
      password: formFields.password,
      attributes: {
        email: formFields.email
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
    displayInfo(`Resending confirmation code for ${formFields.email}`);
    Auth.resendSignUp(formFields.email).catch(displayError);
  }

  function confirmSignUp() {
    displayInfo("Confirming user...");
    Auth.confirmSignUp(formFields.email, formFields.confirmationCode)
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

  function signIn() {
    displayInfo("Signing in...");
    Auth.signIn(formFields.email, formFields.password)
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

  async function startPasswordReset() {
    if (!formFields.email) {
      displayError("Please enter your email first");
      return;
    }
    try {
      await Auth.forgotPassword(formFields.email);
      displayInfo(
        "Please enter your new password and the confirmation code sent to your email address"
      );
      flowState = FlowState.forgotPassword;
    } catch (e) {
      displayError(e);
    }
  }

  function cancelPasswordReset() {
    clearMessages();
    flowState = FlowState.initial;
  }

  async function confirmResetPassword() {
    displayInfo("Resetting password...");
    try {
      await Auth.forgotPasswordSubmit(
        formFields.email,
        formFields.confirmationCode,
        formFields.newPassword
      );
      displayInfo("Please log in with your new password");
      flowState = FlowState.initial;
    } catch (e) {
      displayError(e);
    }
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
        <form on:submit|preventDefault={signIn}>
          <div class="form-group">
            <label>Email Address</label>
            <input
              type="email"
              class="form-control"
              bind:value={formFields.email} />
          </div>
          <div class="form-group">
            <label>Password</label>
            <input
              type="button"
              class="btn btn-clear btn-sm text-muted"
              on:click={startPasswordReset}
              value="(Forgot password?)" />
            <input
              type="password"
              class="form-control"
              bind:value={formFields.password} />
          </div>
          <input
            type="button"
            class="btn btn-secondary"
            on:click={signUp}
            value="Register" />
          <input type="submit" class="btn btn-primary" value="Log In" />
        </form>
      {:else if flowState === FlowState.signUpConfirmation}
        <form on:submit|preventDefault={confirmSignUp}>
          <div class="form-group">
            <label>Confirmation Code</label>
            <input
              type="text"
              class="form-control"
              bind:value={formFields.confirmationCode} />
          </div>
          <input
            class="btn btn-secondary"
            on:click={resendCode}
            value="Resend Confirmation Code" />
          <input type="submit" class="btn btn-primary" value="Confirm" />
        </form>
      {:else if flowState === FlowState.forgotPassword}
        <form on:submit|preventDefault={confirmResetPassword}>
          <div class="form-group">
            <label>New Password</label>
            <input
              type="password"
              class="form-control"
              bind:value={formFields.newPassword} />
          </div>
          <div class="form-group">
            <label>Confirmation Code</label>
            <input
              type="text"
              class="form-control"
              bind:value={formFields.confirmationCode} />
          </div>
          <input
            class="btn btn-secondary"
            on:click={cancelPasswordReset}
            value="Cancel Reset" />
          <input
            class="btn btn-secondary"
            on:click={startPasswordReset}
            value="Resend Confirmation Code" />
          <input type="submit" class="btn btn-primary" value="Confirm" />
        </form>
      {/if}
    </div>
  </div>
</div>
