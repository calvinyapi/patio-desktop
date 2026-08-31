<script lang="ts">
  import { database } from '../../../wailsjs/go/models';
  import CameraViewWindow from './CameraViewWindow.svelte';

  let { camera }: { camera: database.Camera } = $props();

  let showView: boolean = $state(false);
</script>

<button class="camera-card" onclick={() => (showView = true)}>
  <span class="status" class:active={camera.is_active}></span>
  <h5>{camera.name}</h5>
  <p>{camera.rtsp_url}</p>
</button>

<CameraViewWindow camera={camera} bind:visible={showView} />



<style>
.camera-card {
  box-sizing: border-box;
  border-radius: var(--radius-md);
  overflow: hidden; /* pour que le contenu (ex: image) respecte aussi les coins arrondis */
  border: 1px solid #eee;
  background: #fff;
  box-shadow: var(--shadow-sm);
  transition: box-shadow 0.2s ease, transform 0.2s ease;
  padding: 1.25rem;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 6px;
  width: 100%;
  font: inherit;
  text-align: left;
  cursor: pointer;
}

.camera-card:hover {
  box-shadow: 0 12px 28px rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}

.status {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: #bbb;
}

.status.active {
  background: #2ecc71;
}

h5 {
  margin: 0;
  font-size: 1rem;
  font-weight: 600;
}

p {
  margin: 0;
  font-size: 0.75rem;
  font-family: ui-monospace, SFMono-Regular, "Roboto Mono", Consolas, monospace;
  color: #777;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

</style>