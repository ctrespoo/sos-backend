<script lang="ts">
  import { onMount } from "svelte";
  import { firebaseAuth } from "../../lib/authFirebase";

  onMount(() => {
    firebaseAuth.auth().onAuthStateChanged(async (user) => {
      console.log("local: ", window.location.pathname);
      if (
        !user &&
        !window.location.pathname.startsWith("/auth/login") &&
        !window.location.pathname.startsWith("/auth/criar")
      ) {
        console.log("user is null");
        localStorage.removeItem("token");
        window.location.href = "/auth/login";
        return;
      }
      if (
        user &&
        (window.location.pathname.startsWith("/auth/login") ||
          window.location.pathname.startsWith("/auth/criar"))
      ) {
        window.location.href = "/";
      }
      await user
        ?.getIdToken()
        .then((token) => {
          localStorage.setItem("token", token);
        })
        .catch((error) => {
          console.log(error);
        });

      console.log("onAuthStateChanged: ", user?.displayName);
    });
  });
</script>

<div>
  <slot />
</div>
