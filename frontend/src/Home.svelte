<script lang="ts">
  import { onMount } from "svelte";
  import { CountCameras, CountEvents } from "../wailsjs/go/main/App";

  let cameraCount: number | null = $state(null);
  let eventCount: number | null = $state(null);

  onMount(async () => {
    cameraCount = await CountCameras();
    eventCount = await CountEvents();
  });
</script>

<h1>Welcome to Patio</h1>

<div class="stats">
  <div class="stat-card">
    <div class="stat-header">
      <svg
        class="stat-icon"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
        aria-hidden="true"
      >
        <path
          d="M23 19a2 2 0 0 1-2 2H3a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h4l2-3h6l2 3h4a2 2 0 0 1 2 2z"
        />
        <circle cx="12" cy="13" r="4" />
      </svg>
      <span class="stat-label">Cameras</span>
    </div>
    <span class="stat-value">{cameraCount ?? "—"}</span>
  </div>

  <div class="stat-card">
    <div class="stat-header">
      <svg
        class="stat-icon"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
        aria-hidden="true"
      >
        <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9" />
        <path d="M13.73 21a2 2 0 0 1-3.46 0" />
      </svg>
      <span class="stat-label">Events</span>
    </div>
    <span class="stat-value">{eventCount ?? "—"}</span>
  </div>
</div>

<style>
  .stats {
    display: flex;
    gap: 16px;
    margin-top: 1.75rem;
  }

  .stat-card {
    box-sizing: border-box;
    width: 180px;
    background: var(--color-surface);
    border: 1px solid var(--color-border-subtle);
    border-radius: var(--radius-md);
    box-shadow: var(--shadow-sm);
    padding: 1.25rem;
    display: flex;
    flex-direction: column;
    gap: 0.6rem;
  }

  .stat-header {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .stat-icon {
    width: 16px;
    height: 16px;
    flex-shrink: 0;
    color: var(--color-text-secondary);
  }

  .stat-label {
    font-size: 0.75rem;
    font-weight: 500;
    color: var(--color-text-secondary);
  }

  .stat-value {
    font-size: 2.25rem;
    font-weight: 700;
    letter-spacing: -0.02em;
    color: var(--color-text);
  }
</style>
