<script lang="ts">
  import ScanCamerasDialog from "./components/camera/ScanCamerasDialog.svelte";
  import CameraList from "./components/camera/CameraList.svelte";
  import { GetCameras } from "../wailsjs/go/main/App";
  import { onMount } from "svelte";
  import { database } from "../wailsjs/go/models";

  let visible: boolean = $state(false);
  let cameras: database.Camera[] = $state([]);

  onMount(async () => {
    console.log("the component has mounted");
    cameras = await GetCameras();
  });
</script>

<button
  class="btn-add-camera"
  onclick={() => {
    visible = !visible;
  }}>ajouter une caméra</button
>
<ScanCamerasDialog bind:visible />
<CameraList {cameras} />

<style>
  .btn-add-camera {
    all: unset;
    box-sizing: border-box;
    background: #1a1a1a;
    color: #fff;
    font-size: 0.85rem;
    font-weight: 500;
    padding: 0.6rem 1rem;
    border-radius: var(--radius-sm);
    cursor: pointer;
    transition: background-color 0.15s ease;
    margin-bottom: 1.5rem;
  }

  .btn-add-camera:hover {
    background: #333;
  }
</style>
