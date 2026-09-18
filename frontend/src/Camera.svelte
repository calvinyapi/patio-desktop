<script lang="ts">
  import ScanCamerasDialog from "./components/camera/ScanCamerasDialog.svelte";
  import CameraList from "./components/camera/CameraList.svelte";
  import { GetCameras } from "../wailsjs/go/main/App";
  import { onMount } from "svelte";
  import { database } from "../wailsjs/go/models";

  let visible: boolean = $state(false);
  let cameras: database.Camera[] = $state([]);

  let alertMessage: string | null = $state(null);
  let alertTimeout: ReturnType<typeof setTimeout> | undefined;

  onMount(async () => {
    console.log("the component has mounted");
    cameras = await GetCameras();
  });

  function showAlert(message: string) {
    alertMessage = message;
    clearTimeout(alertTimeout);
    alertTimeout = setTimeout(() => (alertMessage = null), 3000);
  }

  function handleCameraDeleted(id: number) {
    cameras = cameras.filter((c) => c.id !== id);
    showAlert("Caméra supprimée avec succès");
  }

  async function handleCameraAdded() {
    visible = false; // referme la fenêtre d'ajout
    cameras = await GetCameras(); // fait apparaître la nouvelle caméra dans la liste
    showAlert("Caméra ajoutée avec succès");
  }
</script>

{#if alertMessage}
  <div class="alert-success">{alertMessage}</div>
{/if}

<button
  class="btn-add-camera"
  onclick={() => {
    visible = !visible;
  }}>ajouter une caméra</button
>
<ScanCamerasDialog bind:visible onadded={handleCameraAdded} />
<CameraList {cameras} ondeleted={handleCameraDeleted} />

<style>
  .btn-add-camera {
    all: unset;
    box-sizing: border-box;
    background: var(--color-accent);
    color: var(--color-accent-contrast);
    font-size: 0.85rem;
    font-weight: 500;
    padding: 0.6rem 1rem;
    border-radius: var(--radius-sm);
    cursor: pointer;
    transition: background-color 0.15s ease;
    margin-bottom: 1.5rem;
  }

  .btn-add-camera:hover {
    background: var(--color-accent-hover);
  }

  .alert-success {
    position: fixed;
    top: 1.5rem;
    left: 50%;
    transform: translateX(-50%);
    z-index: 200; /* au-dessus des fenêtres modales (z-index: 100) */
    padding: 0.65rem 1.1rem;
    /* même vert que le point de statut "active" des caméras */
    background: rgba(46, 204, 113, 0.15);
    border: 1px solid rgba(46, 204, 113, 0.4);
    color: #22c55e; /* reste lisible clair et sombre */
    font-size: 0.85rem;
    font-weight: 500;
    border-radius: var(--radius-sm);
    box-shadow: var(--shadow-sm);
  }
</style>
