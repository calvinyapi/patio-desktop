<script lang="ts">
  import { database } from "../../../wailsjs/go/models";
  import CameraCard from "./CameraCard.svelte";
  let {
    cameras,
    ondeleted,
  }: { cameras: database.Camera[]; ondeleted?: (id: number) => void } = $props();
</script>

<div class="camera-list">
  {#each cameras as camera (camera.id)}
    <CameraCard {camera} {ondeleted} />
  {/each}
</div>

<style>
  .camera-list {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 16px;
  }
</style>
