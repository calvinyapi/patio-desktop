<script lang="ts">
  import Sidebar from "./components/Sidebar.svelte";
  import Camera from "./Camera.svelte";

  let currentPage: string = $state("home");

  const pagechange = (page: string) => {
    currentPage = page;
  };
</script>

<div class="layout">
  <Sidebar {pagechange} />

  <main>
    {#if currentPage === "home"}
      <h1>Welcome to Patio</h1>
    {:else if currentPage === "cameras"}
      <Camera />
    {:else if currentPage === "events"}
      <h1>Events</h1>
    {:else if currentPage === "settings"}
      <h1>Settings</h1>
    {:else}
      <h1>Page not found</h1>
    {/if}
  </main>
</div>

<style>
  .layout {
    display: flex;
    min-height: 100vh;
  }
  main {
    flex: 1;
    padding: 2.5rem;
  }
</style>
