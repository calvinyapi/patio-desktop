<script lang="ts">
  import { onMount } from "svelte";
  import { database } from "../../../wailsjs/go/models";
  import {
    AddZone,
    DeleteCamera,
    GetZonesByCamera,
    UpdateZone,
  } from "../../../wailsjs/go/main/App";

  let {
    camera,
    visible = $bindable(),
    ondeleted,
  }: {
    camera: database.Camera;
    visible: boolean;
    ondeleted?: (id: number) => void;
  } = $props();

  let critique = $state(true); // zone critique par défaut (cf. is_critical DEFAULT true)

  // window.confirm() n'est pas fiable dans la webview Wails (souvent aucune
  // boîte de dialogue ne s'affiche et l'appel renvoie false) -> confirmation
  // repliée directement dans l'entête au lieu d'un confirm() natif.
  let confirmingDelete = $state(false);

  // id de la zone existante pour cette caméra, si l'app en a déjà une :
  // null tant qu'on n'a rien chargé/enregistré -> save() fera un AddZone,
  // sinon un UpdateZone. Cette fenêtre ne gère qu'une seule zone par caméra.
  let zoneId: number | null = $state(null);

  let containerEl: HTMLDivElement;
  let containerWidth: number = $state(0);
  let containerHeight: number = $state(0);

  let x1: number = $state(0);
  let y1: number = $state(0);
  let x2: number = $state(0);
  let y2: number = $state(0);
  let initialized = false;

  let dragging: "tl" | "br" | null = $state(null);

  $effect(() => {
    if (!initialized && containerWidth > 0 && containerHeight > 0) {
      x1 = containerWidth * 0.25;
      y1 = containerHeight * 0.25;
      x2 = containerWidth * 0.75;
      y2 = containerHeight * 0.75;
      initialized = true;
    }
  });

  // Charge la zone existante de cette caméra (s'il y en a une) pour que
  // "Enregistrer" mette à jour au lieu de dupliquer. Ce composant n'est monté
  // qu'une fois par CameraCard (visible bascule juste son affichage).
  onMount(async () => {
    const zones = await GetZonesByCamera(camera.id);
    const existing = zones[0];
    if (existing) {
      zoneId = existing.id;
      x1 = existing.x1;
      y1 = existing.y1;
      x2 = existing.x2;
      y2 = existing.y2;
      critique = existing.is_critical;
      initialized = true; // empêche le $effect de remettre le rectangle par défaut
    }
  });

  function updateFromPointer(e: PointerEvent, corner: "tl" | "br") {
    if (!containerEl) return;
    const rect = containerEl.getBoundingClientRect();
    const px = Math.min(Math.max(e.clientX - rect.left, 0), rect.width);
    const py = Math.min(Math.max(e.clientY - rect.top, 0), rect.height);
    if (corner === "tl") {
      x1 = px;
      y1 = py;
    } else {
      x2 = px;
      y2 = py;
    }
  }

  function onHandlePointerDown(corner: "tl" | "br") {
    return (e: PointerEvent) => {
      e.preventDefault();
      dragging = corner;
      // Pointer capture keeps every move/up event bound to this handle even
      // when the cursor passes over the iframe (a separate document that
      // would otherwise steal the events and freeze the drag mid-way).
      (e.currentTarget as HTMLElement).setPointerCapture(e.pointerId);
      updateFromPointer(e, corner);
    };
  }

  function onHandlePointerMove(corner: "tl" | "br") {
    return (e: PointerEvent) => {
      if (dragging !== corner) return;
      updateFromPointer(e, corner);
    };
  }

  function onHandlePointerUp(e: PointerEvent) {
    dragging = null;
    (e.currentTarget as HTMLElement).releasePointerCapture(e.pointerId);
  }

  function close() {
    visible = false;
    confirmingDelete = false; // repart propre à la prochaine ouverture
  }

  async function deleteCamera() {
    await DeleteCamera(camera.id);
    ondeleted?.(camera.id);
    close();
  }

  async function save() {
    const zone = database.Zone.createFrom({
      id: zoneId ?? 0,
      camera_id: camera.id,
      name: "",
      x1: Math.round(Math.min(x1, x2)),
      y1: Math.round(Math.min(y1, y2)),
      x2: Math.round(Math.max(x1, x2)),
      y2: Math.round(Math.max(y1, y2)),
      threshold_seconds: 5,
      is_critical: critique,
    });

    if (zoneId !== null) {
      await UpdateZone(zone);
    } else {
      zoneId = await AddZone(zone); // les prochains "Enregistrer" mettront à jour
    }
    close();
  }

  function onKeydown(e: KeyboardEvent) {
    if (e.key === "Escape") close();
  }
</script>

<svelte:window onkeydown={visible ? onKeydown : undefined} />

<div class="overlay" class:hidden={!visible}>
  <div class="panel">
    <button class="btn-close" onclick={close} aria-label="Fermer">&times;</button>

    <div class="header">
      <h2>{camera.name}</h2>
      {#if confirmingDelete}
        <div class="confirm-delete">
          <span>Supprimer cette caméra ?</span>
          <button class="btn-cancel" onclick={() => (confirmingDelete = false)}>Annuler</button>
          <button class="btn-confirm" onclick={deleteCamera}>Oui, supprimer</button>
        </div>
      {:else}
        <button class="btn-delete" onclick={() => (confirmingDelete = true)}>
          Supprimer la caméra
        </button>
      {/if}
    </div>

    <div class="content">
      <div
        class="video-wrap"
        bind:this={containerEl}
        bind:clientWidth={containerWidth}
        bind:clientHeight={containerHeight}
      >
        <iframe
          title={camera.name}
          src={`http://localhost:1984/stream.html?src=camera_${camera.id}`}
          frameborder="0"
          allow="autoplay; fullscreen"
        ></iframe>

        <div class="zone-overlay">
          <div
            class="zone-rect"
            style={`left:${Math.min(x1, x2)}px; top:${Math.min(y1, y2)}px; width:${Math.abs(x2 - x1)}px; height:${Math.abs(y2 - y1)}px;`}
          ></div>
          <button
            class="handle"
            style={`left:${x1}px; top:${y1}px;`}
            onpointerdown={onHandlePointerDown("tl")}
            onpointermove={onHandlePointerMove("tl")}
            onpointerup={onHandlePointerUp}
            aria-label="Coin haut-gauche de la zone"
          ></button>
          <button
            class="handle"
            style={`left:${x2}px; top:${y2}px;`}
            onpointerdown={onHandlePointerDown("br")}
            onpointermove={onHandlePointerMove("br")}
            onpointerup={onHandlePointerUp}
            aria-label="Coin bas-droit de la zone"
          ></button>
        </div>
      </div>

      <div class="side">
        <div class="instructions">
          <span class="instructions-icon">i</span>
          <p>
            Faites glisser les deux coins du rectangle sur l'image pour
            ajuster la zone de détection.
          </p>
        </div>

        <div class="zone-inputs">
          <div class="coords">
            <label>
              x1
              <input type="number" bind:value={x1} />
            </label>
            <label>
              y1
              <input type="number" bind:value={y1} />
            </label>
            <label>
              x2
              <input type="number" bind:value={x2} />
            </label>
            <label>
              y2
              <input type="number" bind:value={y2} />
            </label>
          </div>

          <label>
            critique
            <select bind:value={critique}>
              <option value={false}>Non</option>
              <option value={true}>Oui</option>
            </select>
          </label>
        </div>

        <button class="btn-save" onclick={save}>Enregistrer</button>
      </div>
    </div>
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
    background: var(--color-surface);
    border-radius: var(--radius-lg);
    padding: 2rem;
    width: 920px;
    max-width: calc(100vw - 2rem);
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
    color: var(--color-text-muted);
    cursor: pointer;
    transition:
      background-color 0.1s ease,
      color 0.1s ease;
  }

  .btn-close:hover {
    background: var(--color-hover);
    color: var(--color-text);
  }

  .header {
    display: flex;
    align-items: baseline;
    justify-content: space-between;
    gap: 1rem;
    margin-bottom: 1rem;
    padding-right: 1.5rem; /* laisse la place au bouton "Fermer" */
  }

  .header h2 {
    margin: 0;
  }

  .btn-delete {
    all: unset;
    box-sizing: border-box;
    flex-shrink: 0;
    font-size: 0.75rem;
    color: #b3413a;
    padding: 0.35rem 0.6rem;
    border-radius: var(--radius-sm);
    cursor: pointer;
    transition:
      background-color 0.12s ease,
      color 0.12s ease;
  }

  .btn-delete:hover {
    /* lavis translucide plutôt qu'un fond clair fixe : reste lisible en mode nuit */
    background: rgba(179, 65, 58, 0.12);
  }

  .confirm-delete {
    flex-shrink: 0;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.75rem;
    color: var(--color-text-secondary);
  }

  .btn-cancel {
    all: unset;
    box-sizing: border-box;
    font-size: 0.75rem;
    color: var(--color-text-secondary);
    padding: 0.35rem 0.6rem;
    border-radius: var(--radius-sm);
    cursor: pointer;
    transition: background-color 0.12s ease;
  }

  .btn-cancel:hover {
    background: var(--color-bg-subtle);
  }

  .btn-confirm {
    all: unset;
    box-sizing: border-box;
    font-size: 0.75rem;
    font-weight: 500;
    color: #fff;
    background: #b3413a;
    padding: 0.35rem 0.6rem;
    border-radius: var(--radius-sm);
    cursor: pointer;
    transition: background-color 0.12s ease;
  }

  .btn-confirm:hover {
    background: #8f2d2d;
  }

  .content {
    display: flex;
    gap: 1.25rem;
    align-items: flex-start;
  }

  .video-wrap {
    position: relative;
    flex: 1;
    min-width: 0;
  }

  iframe {
    display: block;
    width: 100%;
    height: 440px;
    border-radius: var(--radius-md);
    border: 1px solid var(--color-border);
  }

  .zone-overlay {
    position: absolute;
    inset: 0;
    pointer-events: none;
  }

  .zone-rect {
    position: absolute;
    box-sizing: border-box;
    border: 2px solid #fff;
    box-shadow:
      0 0 0 1px rgba(0, 0, 0, 0.45),
      inset 0 0 0 1px rgba(0, 0, 0, 0.45);
    background: rgba(255, 255, 255, 0.12);
    border-radius: 2px;
  }

  .handle {
    all: unset;
    position: absolute;
    box-sizing: border-box;
    width: 14px;
    height: 14px;
    border-radius: 50%;
    background: #fff;
    border: 2px solid var(--color-text);
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.35);
    transform: translate(-50%, -50%);
    pointer-events: auto;
    touch-action: none;
    cursor: grab;
  }

  .handle:active {
    cursor: grabbing;
  }

  .side {
    flex-shrink: 0;
    width: 196px;
    min-height: 440px; /* aligne le bas de la colonne sur celui de la vidéo */
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .instructions {
    display: flex;
    align-items: flex-start;
    gap: 0.5rem;
    padding: 0.7rem 0.75rem;
    background: var(--color-bg-subtle);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
  }

  .instructions-icon {
    flex-shrink: 0;
    width: 16px;
    height: 16px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    background: var(--color-text-muted);
    color: #fff;
    font-size: 0.65rem;
    font-style: italic;
    font-family: Georgia, serif;
  }

  .instructions p {
    margin: 0;
    font-size: 0.75rem;
    line-height: 1.4;
    color: var(--color-text-secondary);
  }

  .zone-inputs {
    display: flex;
    flex-direction: column;
    gap: 0.85rem;
  }

  .coords {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 0.55rem;
  }

  .zone-inputs label {
    display: flex;
    flex-direction: column;
    gap: 0.35rem;
    font-size: 0.75rem;
    color: var(--color-text-secondary);
  }

  .zone-inputs input,
  .zone-inputs select {
    box-sizing: border-box;
    width: 100%;
    height: 34px;
    font-family: inherit;
    font-size: 0.9rem;
    color: var(--color-text);
    padding: 0 0.6rem;
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    transition: border-color 0.12s ease;
  }

  .zone-inputs input:focus,
  .zone-inputs select:focus {
    outline: none;
    border-color: var(--color-text);
  }

  .zone-inputs input::-webkit-outer-spin-button,
  .zone-inputs input::-webkit-inner-spin-button {
    -webkit-appearance: none;
    margin: 0;
  }

  .zone-inputs select {
    appearance: none;
    -webkit-appearance: none;
    cursor: pointer;
    padding-right: 1.6rem;
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='10' height='6' viewBox='0 0 10 6'%3E%3Cpath d='M1 1l4 4 4-4' fill='none' stroke='%23999' stroke-width='1.5' stroke-linecap='round' stroke-linejoin='round'/%3E%3C/svg%3E");
    background-repeat: no-repeat;
    background-position: right 0.65rem center;
  }

  .btn-save {
    all: unset;
    box-sizing: border-box;
    display: block;
    width: 100%;
    margin-top: auto; /* ancre le bouton en bas de la colonne */
    text-align: center;
    background: var(--color-accent);
    color: var(--color-accent-contrast);
    font-size: 0.85rem;
    font-weight: 500;
    padding: 0.6rem 1rem;
    border-radius: var(--radius-sm);
    cursor: pointer;
    transition: background-color 0.15s ease;
  }

  .btn-save:hover {
    background: var(--color-accent-hover);
  }
</style>
