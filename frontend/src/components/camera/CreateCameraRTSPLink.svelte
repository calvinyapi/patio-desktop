<script lang="ts">
  let camera_name: string = $state("");
  let rtsp_url: string = $state("");
  import { AddCamera } from "../../../wailsjs/go/main/App";

  let { onadded }: { onadded?: () => void } = $props();

  async function CreateCameraViaRTSPLink() {
    try {
      await AddCamera(camera_name, rtsp_url);
      onadded?.();
    } catch (e) {
      console.error("Erreur lors de la création de la caméra :", e);
    }
  }
</script>

<form
  onsubmit={(e) => {
    e.preventDefault();
    CreateCameraViaRTSPLink();
  }}
>
  <div class="form-group">
    <label for="camera-name">Nom de la caméra</label>
    <input type="text" id="camera-name" bind:value={camera_name} />
  </div>

  <div class="form-group">
    <label for="camera-endpoint">Endpoint</label>
    <input
      type="text"
      id="camera-endpoint"
      bind:value={rtsp_url}
      placeholder="rtsp://utilisateur:motdepasse@192.168.1.50:554/stream"
    />
  </div>

  <button type="submit" class="btn-submit">Créer</button>
</form>

<style>
  .form-group {
    display: flex;
    flex-direction: column;
    gap: 0.35rem;
    margin-bottom: 1.1rem;
  }

  .form-group label {
    font-size: 0.8rem;
    color: var(--color-text-secondary);
  }

  .form-group input {
    box-sizing: border-box;
    font-size: 0.9rem;
    padding: 0.6rem 0.7rem;
    background: var(--color-surface);
    color: var(--color-text);
    border: 1px solid var(--color-border-strong);
    border-radius: var(--radius-sm);
  }

  .form-group input:focus {
    outline: none;
    border-color: var(--color-text);
  }

  .form-group input::placeholder {
    color: var(--color-text-muted);
  }

  .btn-submit {
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
  }

  .btn-submit:hover {
    background: var(--color-accent-hover);
  }
</style>
