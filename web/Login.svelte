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

  onMount(() => {
    Auth.currentAuthenticatedUser()
      .then(signedIn)
      .catch(err => {
        console.log("Not signed in:", err);
        flowState = FlowState.initial;
      });
  });

  function signedIn() {
    dispatch("signedIn");
  }

  function signUp({ email, password }) {
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
          flowState = FlowState.signUpConfirmation;
          return;
        }
        signedIn();
      })
      .catch(console.error);
  }

  function resendCode() {
    console.log("resending:", _username);
    Auth.resendSignUp(_username).catch(console.error);
  }

  function confirmSignUp({ detail: data }) {
    console.log("confirming:", _username, data);
    Auth.confirmSignUp(_username, data.code).then(result => {
      if (result === "SUCCESS") {
        console.log("Signed in");
        signedIn();
      }
    });
  }

  function signIn({ detail: data }) {
    _username = data.email;
    Auth.signIn(data.email, data.password)
      .then(signedIn)
      .catch(err => {
        if (err.code === "UserNotConfirmedException") {
          flowState = FlowState.signUpConfirmation;
        }
      });
  }
</script>

<div class="container">
  <div class="card">
    <div class="card-body">
      <h5 class="card-title">Log in</h5>
      {#if flowState === FlowState.initial}
        <Form
          bind:this={mainForm}
          formFields={initialFormFields}
          on:submit={signIn}>
          <button
            class="btn btn-secondary"
            on:click|preventDefault={() => signUp(mainForm.getFormData().details)}>
            Register
          </button>
          <button type="submit" class="btn btn-primary">Log In</button>
        </Form>
      {:else if flowState === FlowState.signUpConfirmation}
        <Form
          bind:this={mainForm}
          formFields={signUpConfirmationFormFields}
          on:submit={confirmSignUp}>
          <button
            class="btn btn-secondary"
            on:click|preventDefault={() => resendCode(mainForm.getFormData().details)}>
            Resend Confirmation Code
          </button>
          <button type="submit" class="btn btn-primary">Confirm</button>
        </Form>
      {/if}
    </div>
  </div>
</div>
