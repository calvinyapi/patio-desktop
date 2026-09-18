<script lang="ts">
  import Sidebar from "./components/Sidebar.svelte";
  import Camera from "./Camera.svelte";
  import Home from "./Home.svelte";
  import Settings from "./Settings.svelte";

  let currentPage: string = $state("home");

  const pagechange = (page: string) => {
    currentPage = page;
  };

  const pageLabels: Record<string, string> = {
    home: "Dashboard",
    cameras: "Cameras",
    events: "Events",
    settings: "Settings",
  };
</script>

<div class="layout">
  <Sidebar {pagechange} />

  <main>
    <nav class="breadcrumb" aria-label="Fil d'Ariane">
      <button class="crumb-root" onclick={() => pagechange("home")}>Patio</button>
      <span class="crumb-sep">/</span>
      <span class="crumb-current">{pageLabels[currentPage] ?? "Page introuvable"}</span>
    </nav>

    {#if currentPage === "home"}
      <Home />
    {:else if currentPage === "cameras"}
      <Camera />
    {:else if currentPage === "events"}
      <h1>Events</h1>
    {:else if currentPage === "settings"}
      <Settings />
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

  .breadcrumb {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    margin-bottom: 1.5rem;
    font-size: 0.8rem;
  }

  .crumb-root {
    all: unset;
    color: var(--color-text-secondary);
    cursor: pointer;
    transition: color 0.12s ease;
  }

  .crumb-root:hover {
    color: var(--color-text);
  }

  .crumb-sep {
    color: var(--color-text-muted);
  }

  .crumb-current {
    color: var(--color-text);
    font-weight: 500;
  }
</style>
