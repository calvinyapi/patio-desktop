<script lang="ts">
  import { database } from "../../../wailsjs/go/models";

  let {
    camera,
    visible = $bindable(),
  }: { camera: database.Camera; visible: boolean } = $props();

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
  }

  function onKeydown(e: KeyboardEvent) {
    if (e.key === "Escape") close();
  }
</script>

<svelte:window onkeydown={visible ? onKeydown : undefined} />

<div class="overlay" class:hidden={!visible}>
  <div class="panel">
    <button class="btn-close" onclick={close} aria-label="Fermer">&times;</button>

    <h2>{camera.name}</h2>

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
    background: #fff;
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

  h2 {
    margin-bottom: 1rem;
    padding-right: 1.5rem;
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
    width: 180px;
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
    gap: 0.6rem;
  }

  .zone-inputs label {
    display: flex;
    flex-direction: column;
    gap: 0.3rem;
    font-size: 0.75rem;
    color: var(--color-text-secondary);
  }

  .zone-inputs input {
    box-sizing: border-box;
    width: 100%;
    font-size: 0.9rem;
    padding: 0.4rem 0.5rem;
    border: 2px solid var(--color-border);
    border-radius: var(--radius-sm);
  }
</style>
