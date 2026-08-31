<script lang="ts">
  import { models } from "../../../wailsjs/go/models";
  import { ScanNetwork } from "../../../wailsjs/go/main/App";
  import AddCameraForm from "./CreateCameraAutoDialog.svelte";
  import CreateCameraRTSPLink from "./CreateCameraRTSPLink.svelte";

  let { visible = $bindable() } = $props();
  let scanning: boolean = $state(true);
  let discoveredCameras: models.DiscoveredCamera[] = $state([]);
  let selectedCamera: models.DiscoveredCamera | null = $state(null);
  let showManualForm: boolean = $state(false);

  async function ScanNetworkCameras() {
    try {
      scanning = true;
      discoveredCameras = await ScanNetwork();
      console.log("Caméras découvertes :", discoveredCameras);
    } catch (error) {
      console.error("Erreur lors de la découverte des caméras :", error);
    } finally {
      scanning = false;
    }
  }

  $effect(() => {
    console.log("effect déclenché, visible =", visible);
    if (visible) {
      ScanNetworkCameras();
    }
  });
</script>

<div class="overlay" class:hidden={!visible}>
  <div class="panel">
    <button
      class="btn-close"
      onclick={() => {
        visible = false;
      }}
      aria-label="Fermer">&times;</button
    >
    {#if selectedCamera}
      <AddCameraForm bind:selectedCamera />
    {:else if showManualForm}
      <button
        class="btn-back"
        onclick={() => {
          showManualForm = false;
        }}>&larr; Retour</button
      >
      <h2>Ajouter via lien RTSP</h2>
      <CreateCameraRTSPLink />
    {:else}
      <h2>Scan du réseau</h2>
      <div class="camera-shelf">
        {#if scanning}
          <p class="empty">Recherche de caméras en cours...</p>
        {:else}
          {#each discoveredCameras as camera}
            <button
              class="camera-item"
              onclick={() => {
                selectedCamera = camera;
              }}
            >
              <span class="camera-name">{camera.name}</span>
              <span class="camera-endpoint">{camera.endpoint}</span>
            </button>
          {:else}
            <p class="empty">Aucune caméra découverte.</p>
          {/each}
        {/if}
      </div>

      <button
        class="btn-manual"
        onclick={() => {
          showManualForm = true;
        }}>Ajouter une caméra via lien RTSP</button
      >
    {/if}
  </div>
</div>

<style>
  .overlay {
    position: fixed;
    inset: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(0, 0, 0, 0.25);
    backdrop-filter: blur(8px);
    -webkit-backdrop-filter: blur(8px);
    z-index: 100;
  }

  .overlay.hidden {
    display: none;
  }

  .panel {
    position: relative;
    background: #fff;
    border-radius: var(--radius-lg);
    padding: 2rem;
    min-width: 360px;
    box-shadow: var(--shadow-md);
  }

  .btn-close {
    all: unset;
    position: absolute;
    top: 0.75rem;
    right: 0.75rem;
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 4px;
    font-size: 1.1rem;
    line-height: 1;
    color: #888;
    cursor: pointer;
    transition:
      background-color 0.1s ease,
      color 0.1s ease;
  }

  .btn-close:hover {
    background: #f0f0f0;
    color: #222;
  }

  .btn-back {
    all: unset;
    font-size: 0.8rem;
    color: #888;
    cursor: pointer;
    margin-bottom: 0.75rem;
  }

  .btn-back:hover {
    color: #222;
  }

  .camera-shelf {
    margin-top: 1rem;
    max-height: 280px;
    overflow-y: auto;
    border: 1px solid #eee;
    border-radius: var(--radius-md);
  }

  .camera-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    padding: 0.6rem 0.85rem;
    border-bottom: 1px solid #eee;
    transition: background-color 0.1s ease;
  }

  .camera-item:last-child {
    border-bottom: none;
  }

  .camera-item:hover {
    background: #fafafa;
  }

  .camera-name {
    font-size: 0.85rem;
    color: #333;
    white-space: nowrap;
  }

  .camera-endpoint {
    font-size: 0.75rem;
    font-family:
      ui-monospace, SFMono-Regular, "Roboto Mono", Consolas, monospace;
    color: #999;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .empty {
    padding: 0.85rem;
    font-size: 0.85rem;
    color: #999;
    text-align: center;
  }

  .btn-manual {
    all: unset;
    box-sizing: border-box;
    display: block;
    width: 100%;
    text-align: center;
    margin-top: 0.75rem;
    padding: 0.5rem 0.9rem;
    border: 1px solid #ddd;
    border-radius: var(--radius-sm);
    font-size: 0.85rem;
    color: #444;
    cursor: pointer;
    transition:
      background-color 0.1s ease,
      border-color 0.1s ease;
  }

  .btn-manual:hover {
    background: #fafafa;
    border-color: #ccc;
  }
</style>
